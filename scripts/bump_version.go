package main

import (
	"fmt"
	"os"
	"regexp"
)

// main 是版本更新脚本的主入口
// 这是一个运维工具脚本，用于在发布新版本时批量修改项目中的所有版本号配置。
func main() {
	// 1. 校验命令行参数：必须提供新的版本号
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run scripts/bump_version.go <new_version>")
		os.Exit(1)
	}
	newVersion := os.Args[1]
	fmt.Printf("Bumping version to %s...\n", newVersion)

	// updateFile 定义一个内部辅助函数，用于通用的文件正则替换操作
	//
	// 参数说明:
	//   - path:    要修改的文件路径
	//   - pattern: 用于定位旧版本号的正则表达式
	//   - repl:    包含新版本号的替换内容
	updateFile := func(path string, pattern string, repl string) {
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", path, err)
			return
		}

		re := regexp.MustCompile(pattern)

		// 检查正则表达式是否匹配，如果没有匹配项则输出警告 (防止正则失效导致未更新)
		if !re.Match(content) {
			fmt.Printf("Warning: No match found in %s for pattern '%s'\n", path, pattern)
			return
		}

		// 执行全局替换
		newContent := re.ReplaceAll(content, []byte(repl))

		// 将更新后的内容写回文件，保持 0644 权限
		if err := os.WriteFile(path, newContent, 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", path, err)
			return
		}
		fmt.Printf("Updated %s\n", path)
	}

	// ---------------------------------------------------------
	// 开始依次更新各平台的配置文件
	// ---------------------------------------------------------

	// 1. 更新后端配置文件 (build/config.yml)
	// 目标格式: version: "x.x.x"
	updateFile("build/config.yml", `version: "\d+\.\d+\.\d+"`, fmt.Sprintf(`version: "%s"`, newVersion))

	// 2. 更新前端项目文件 (frontend/package.json)
	// 目标格式: "version": "x.x.x"
	updateFile("frontend/package.json", `"version": "\d+\.\d+\.\d+"`, fmt.Sprintf(`"version": "%s"`, newVersion))

	// 3. 更新 macOS 应用元数据 (Info.plist & Info.dev.plist)
	// 这里使用正则捕获组 (Capture Groups) 来保留原有的 XML 标签和缩进格式。
	// $1 表示第一个捕获组 (开始标签)，${2} 表示第二个捕获组 (结束标签)。

	// 更新 CFBundleShortVersionString (显示版本号)
	plistPattern := `(<key>CFBundleShortVersionString</key>\s*<string>).*?(</string>)`
	updateFile("build/darwin/Info.plist", plistPattern, fmt.Sprintf("${1}%s${2}", newVersion))

	// 更新 CFBundleVersion (构建版本号)
	plistVersionPattern := `(<key>CFBundleVersion</key>\s*<string>).*?(</string>)`
	updateFile("build/darwin/Info.plist", plistVersionPattern, fmt.Sprintf("${1}%s${2}", newVersion))

	// 【开发环境】同时更新开发环境的 plist 文件
	updateFile("build/darwin/Info.dev.plist", plistPattern, fmt.Sprintf("${1}%s${2}", newVersion))
	updateFile("build/darwin/Info.dev.plist", plistVersionPattern, fmt.Sprintf("${1}%s${2}", newVersion))

	// 4. 更新 Windows 应用元数据 (build/windows/info.json)
	// 需要同时更新 file_version 和 ProductVersion
	updateFile("build/windows/info.json", `"file_version": "\d+\.\d+\.\d+"`, fmt.Sprintf(`"file_version": "%s"`, newVersion))
	updateFile("build/windows/info.json", `"ProductVersion": "\d+\.\d+\.\d+"`, fmt.Sprintf(`"ProductVersion": "%s"`, newVersion))

	// 5. 更新 Linux 打包配置 (build/linux/nfpm/nfpm.yaml)
	// 使用 nfpm 进行打包时的版本配置
	updateFile("build/linux/nfpm/nfpm.yaml", `version: "\d+\.\d+\.\d+"`, fmt.Sprintf(`version: "%s"`, newVersion))

	fmt.Println("🎉 All files updated successfully!")
}
