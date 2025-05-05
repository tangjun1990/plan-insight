package imagex

import (
	"fmt"
	"path/filepath"
)

// Example_drawTriangleOnImage 展示如何使用DrawTriangleOnImage函数
func Example_drawTriangleOnImage() {
	// 示例：加载一个图片，在中间绘制红色三角形，然后保存

	// 获取临时目录
	tempDir := "./" // os.TempDir()

	// 输入图片路径 (实际使用时请替换为实际图片路径)
	inputPath := filepath.Join(tempDir, "test.jpg")

	// 输出图片路径
	outputPath := filepath.Join(tempDir, "output_with_triangle.jpg")

	// 调用函数在图片中间添加红色三角形
	err := DrawTriangleOnImage(inputPath, outputPath)
	if err != nil {
		fmt.Printf("处理图片时出错: %v\n", err)
		return
	}

	fmt.Printf("图片处理成功，已保存到: %s\n", outputPath)

	// 注意：由于这个示例需要实际的图片文件才能运行，
	// 实际使用时它不会作为测试运行，而是作为文档示例

	// Output:
}

// Example_drawRectangleOnImage 展示如何使用DrawRectangleOnImage函数
func Example_drawRectangleOnImage() {
	// 示例：加载一个图片，在指定位置绘制蓝色矩形，然后保存

	// 获取临时目录
	tempDir := "./" // os.TempDir()

	// 输入图片路径 (实际使用时请替换为实际图片路径)
	inputPath := filepath.Join(tempDir, "colorbase.jpg")

	// 输出图片路径
	outputPath := filepath.Join(tempDir, "output_with_rectangle.jpg")

	// 设置矩形参数
	x := 250        // 矩形左上角x坐标
	y := 397        // 矩形左上角y坐标
	width := 95     // 矩形宽度
	height := 93    // 矩形高度
	r := uint8(0)   // 红色分量 (0-255)
	g := uint8(0)   // 绿色分量 (0-255)
	b := uint8(255) // 蓝色分量 (0-255)

	// 调用函数在指定位置添加蓝色矩形
	err := DrawRectangleOnImage(inputPath, x, y, width, height, r, g, b, outputPath)
	if err != nil {
		fmt.Printf("处理图片时出错: %v\n", err)
		return
	}

	fmt.Printf("图片处理成功，已保存到: %s\n", outputPath)

	// 注意：由于这个示例需要实际的图片文件才能运行，
	// 实际使用时它不会作为测试运行，而是作为文档示例

	// Output:
}
