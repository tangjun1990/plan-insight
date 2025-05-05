# 图像处理工具包

这个包提供了一系列图像处理的工具函数。

## 功能

### DrawTriangleOnImage

在图片中间添加红色三角形，并保存到指定路径。

#### 参数

- `imagePath`：输入图片的路径（支持JPG和PNG格式）
- `outputPath`：输出图片的路径，如果为空则覆盖原图片

#### 返回值

- `error`：处理过程中发生的错误

### DrawRectangleOnImage

在图片上绘制指定颜色的矩形，并保存到指定路径。

#### 参数

- `imagePath`：输入图片的路径（支持JPG和PNG格式）
- `x`：矩形左上角的x坐标
- `y`：矩形左上角的y坐标
- `width`：矩形宽度
- `height`：矩形高度
- `r, g, b`：矩形颜色的RGB值
- `outputPath`：输出图片的路径，如果为空则覆盖原图片

#### 返回值

- `error`：处理过程中发生的错误

## 使用示例

### 绘制三角形

```go
package main

import (
    "fmt"
    "git.4321.sh/feige/commonapi/pkg/imagex"
)

func main() {
    // 在图片中间添加红色三角形
    err := imagex.DrawTriangleOnImage("input.jpg", "output.jpg")
    if err != nil {
        fmt.Printf("处理图片失败: %v\n", err)
        return
    }
    
    fmt.Println("图片处理成功!")
}
```

### 绘制矩形

```go
package main

import (
    "fmt"
    "git.4321.sh/feige/commonapi/pkg/imagex"
)

func main() {
    // 在图片上绘制蓝色矩形
    // 参数: 输入图片路径, x坐标, y坐标, 宽度, 高度, RGB颜色值, 输出图片路径
    err := imagex.DrawRectangleOnImage("input.jpg", 100, 100, 200, 150, 0, 0, 255, "output_rectangle.jpg")
    if err != nil {
        fmt.Printf("处理图片失败: %v\n", err)
        return
    }
    
    fmt.Println("图片处理成功!")
}
```

## 注意事项

- 支持的图片格式：JPG、PNG
- 三角形大小为图片较短边的1/4
- 三角形为实心红色
- 矩形绘制超出图片边界会自动裁剪
- 坐标系以图片左上角为原点(0,0) 