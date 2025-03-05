package nft

import (
	"fmt"
	"log"
	"net"
	"space-api/util"

	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

type nftService struct {
	conn  *nftables.Conn
	table *nftables.Table
	chain *nftables.Chain
}

var DefaultNftService = NewNftService("drop_ip_table", "drop_ip_chain")

func NewNftService(tableName, chainName string) *nftService {
	// 创建 nftables 连接
	conn := &nftables.Conn{}

	// 创建 IPv4 过滤表
	table := &nftables.Table{
		Family: nftables.TableFamilyINet,
		Name:   tableName,
	}

	conn.AddTable(table)
	// 创建输入链
	chain := &nftables.Chain{
		Name:     chainName,
		Table:    table,
		Type:     nftables.ChainTypeFilter,
		Hooknum:  nftables.ChainHookInput,
		Priority: nftables.ChainPriorityFilter,
	}
	conn.AddChain(chain)

	// 提交初始配置
	if err := conn.Flush(); err != nil {
		log.Fatalf("Failed to create initial nftables configuration: %v", err)
	}

	return &nftService{
		conn:  conn,
		table: table,
		chain: chain,
	}
}

func (m *nftService) AddBlockIPList(ips []string) error {
	// 解析 IP 地址

	rules := []*nftables.Rule{}
	for _, ip := range ips {
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			return util.CreateBizErr(
				fmt.Sprintf("invalid IP address: %s", ip),
				fmt.Errorf("invalid IP address: %s", ip),
			)
		}

		// 创建封禁规则
		rule := &nftables.Rule{
			Table: m.table,
			Chain: m.chain,
			Exprs: []expr.Any{
				// 匹配源 IP
				&expr.Payload{
					DestRegister: 1,
					Base:         expr.PayloadBaseNetworkHeader,
					Offset:       12, // Source IP offset for IPv4
					Len:          4,  // IPv4 地址长度
				},
				// 匹配网络协议
				&expr.Cmp{
					Op:       expr.CmpOpEq,
					Register: 1,
					Data:     parsedIP.To4(),
				},
				// 丢弃匹配的数据包
				&expr.Verdict{
					Kind: expr.VerdictDrop,
				},
			},
		}
		rules = append(rules, rule)
	}

	for _, rule := range rules {
		// 添加规则
		m.conn.AddRule(rule)
	}

	// 提交更改
	if err := m.conn.Flush(); err != nil {
		return util.CreateBizErr("submit ip rule failed", fmt.Errorf("failed to ban IP %s: %v", ips, err))
	}

	return nil
}

func (m *nftService) RemoveBlockIPList(ips []string) error {
	// 解析 IP 地址
	parsedIPList := []*net.IP{}
	for _, ip := range ips {
		parsedIP := net.ParseIP(ip)
		if parsedIP == nil {
			return util.CreateBizErr(
				fmt.Sprintf("invalid IP address: %s", ip),
				fmt.Errorf("invalid IP address: %s", ip),
			)
		}
		parsedIPList = append(parsedIPList, &parsedIP)
	}

	// 获取所有规则
	rules, err := m.conn.GetRules(m.table, m.chain)
	if err != nil {
		return fmt.Errorf("failed to get rules: %v", err)
	}

	for _, parsedIP := range parsedIPList {
		// 遍历规则，找到匹配的 IP 规则并删除
		for _, rule := range rules {
			for _, e := range rule.Exprs {
				if cmp, ok := e.(*expr.Cmp); ok {
					// 检查是否是 IP 匹配规则
					if len(cmp.Data) == 4 && net.IP(cmp.Data).Equal(*parsedIP) {
						m.conn.DelRule(rule)
						break
					}
				}
			}
		}
	}

	// 提交更改
	if err := m.conn.Flush(); err != nil {
		return fmt.Errorf("failed to unban IP %s: %v", ips, err)
	}

	return nil
}

func (m *nftService) GetBlockList() (resp []string, err error) {
	resp = []string{}
	// 获取所有规则
	rules, err := m.conn.GetRules(m.table, m.chain)
	if err != nil {
		err = util.CreateBizErr("failed get rules", fmt.Errorf("failed to get rules: %v", err))
		return
	}

	// 遍历规则，提取被封禁的 IP
	for _, rule := range rules {
		for _, e := range rule.Exprs {
			if cmp, ok := e.(*expr.Cmp); ok {
				// 检查是否是 IP 匹配规则
				if len(cmp.Data) == 4 {
					resp = append(resp, net.IP(cmp.Data).String())
					break
				}
			}
		}
	}

	return
}
