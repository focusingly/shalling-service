package constants

// ANSI color codes
const (
	RESET      = "\033[0m"  // 重置所有属性
	BOLD       = "\033[1m"  // 加粗
	UNDERLINE  = "\033[4m"  // 下划线
	BLACK      = "\033[30m" // 黑色
	RED        = "\033[31m" // 红色
	GREEN      = "\033[32m" // 绿色
	YELLOW     = "\033[33m" // 黄色
	BLUE       = "\033[34m" // 蓝色
	MAGENTA    = "\033[35m" // 紫色
	CYAN       = "\033[36m" // 青色
	WHITE      = "\033[37m" // 白色
	BG_BLACK   = "\033[40m" // 黑色背景
	BG_RED     = "\033[41m" // 红色背景
	BG_GREEN   = "\033[42m" // 绿色背景
	BG_YELLOW  = "\033[43m" // 黄色背景
	BG_BLUE    = "\033[44m" // 蓝色背景
	BG_MAGENTA = "\033[45m" // 紫色背景
	BG_CYAN    = "\033[46m" // 青色背景
	BG_WHITE   = "\033[47m" // 白色背景
)

// Log colors for different levels
const (
	INFO  = GREEN  // INFO level - 绿色
	ERROR = RED    // ERROR level - 红色
	WARN  = YELLOW // WARN level - 黄色
	DEBUG = CYAN   // DEBUG level - 青色
)
