package aesthetic

import (
	"encoding/json"
	"fmt"
	"image/color"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/tangjun1990/flygo/core/kcfg"
)

type colorbarItem struct {
	Num int    `json:"num"`
	URL string `json:"url"`
}

var globalColorbar = []colorbarItem{
	{
		1,
		"/colorbar/1.png",
	},
	{
		2,
		"/colorbar/2.png",
	},
	{
		3,
		"/colorbar/3.png",
	},
	{
		4,
		"/colorbar/4.png",
	},
	{
		5,
		"/colorbar/5.png",
	},
	{
		6,
		"/colorbar/6.png",
	},
	{
		7,
		"/colorbar/7.png",
	},
	{
		8,
		"/colorbar/8.png",
	},
	{
		9,
		"/colorbar/9.png",
	},
	{
		10,
		"/colorbar/10.png",
	},
	{
		11,
		"/colorbar/11.png",
	},
	{
		12,
		"/colorbar/12.png",
	},
	{
		13,
		"/colorbar/13.png",
	},
	{
		14,
		"/colorbar/14.png",
	},
	{
		15,
		"/colorbar/15.png",
	},
	{
		16,
		"/colorbar/16.png",
	},
	{
		17,
		"/colorbar/17.png",
	},
	{
		18,
		"/colorbar/18.png",
	},
	{
		19,
		"/colorbar/19.png",
	},
	{
		20,
		"/colorbar/20.png",
	},
	{
		21,
		"/colorbar/21.png",
	},
	{
		22,
		"/colorbar/22.png",
	},
	{
		23,
		"/colorbar/23.png",
	},
	{
		24,
		"/colorbar/24.png",
	},
	{
		25,
		"/colorbar/25.png",
	},
	{
		26,
		"/colorbar/26.png",
	},
	{
		27,
		"/colorbar/27.png",
	},
	{
		28,
		"/colorbar/28.png",
	},
	{
		29,
		"/colorbar/29.png",
	},
	{
		30,
		"/colorbar/30.png",
	},
	{
		31,
		"/colorbar/31.png",
	},
	{
		32,
		"/colorbar/32.png",
	},
	{
		33,
		"/colorbar/33.png",
	},
	{
		34,
		"/colorbar/34.png",
	},
	{
		35,
		"/colorbar/35.png",
	},
	{
		36,
		"/colorbar/36.png",
	},
	{
		37,
		"/colorbar/37.png",
	},
	{
		38,
		"/colorbar/38.png",
	},
	{
		39,
		"/colorbar/39.png",
	},
	{
		40,
		"/colorbar/40.png",
	},
	{
		41,
		"/colorbar/41.png",
	},
	{
		42,
		"/colorbar/42.png",
	},
	{
		43,
		"/colorbar/43.png",
	},
	{
		44,
		"/colorbar/44.png",
	},
	{
		45,
		"/colorbar/45.png",
	},
	{
		46,
		"/colorbar/46.png",
	},
	{
		47,
		"/colorbar/47.png",
	},
	{
		48,
		"/colorbar/48.png",
	},
}

type boxItemV3 struct {
	Num          int
	name         string
	StartX       int
	StartY       int
	Words        []string
	Colors       []int
	ColorBars    []int
	ColorComment string
	Tag          []string
	Comment      string
	Attitude     string
}

var globalImageV3 = []string{
	"1-1-1.jpg",
	"1-1-2.jpg",
	"1-1-3.jpg",
	"1-1-4.jpg",
	"1-1-5.jpg",
	"1-1-6.jpg",
	"1-1-7.jpg",
	"1-1-8.jpg",
	"1-1-9.jpg",
	"1-1-10.jpg",
	"1-1-11.jpg",
	"1-1-12.jpg",
	"1-1-13.jpg",
	"1-1-14.jpg",
	"1-1-15.jpg",
	"1-1-16.jpg",

	"1-2-1.jpg",
	"1-2-2.jpg",
	"1-2-3.jpg",
	"1-2-4.jpg",
	"1-2-5.jpg",
	"1-2-6.jpg",
	"1-2-7.jpg",
	"1-2-8.jpg",
	"1-2-9.jpg",
	"1-2-10.jpg",
	"1-2-11.jpg",
	"1-2-12.jpg",
	"1-2-13.jpg",
	"1-2-14.jpg",
	"1-2-15.jpg",
	"1-2-16.jpg",

	"2-1-1.jpg",
	"2-1-2.jpg",
	"2-1-3.jpg",
	"2-1-4.jpg",
	"2-1-5.jpg",
	"2-1-6.jpg",
	"2-1-7.jpg",
	"2-1-8.jpg",
	"2-1-9.jpg",
	"2-1-10.jpg",
	"2-1-11.jpg",
	"2-1-12.jpg",
	"2-1-13.jpg",
	"2-1-14.jpg",
	"2-1-15.jpg",
	"2-1-16.jpg",

	"3-1-1.jpg",
	"3-1-2.jpg",
	"3-1-3.jpg",
	"3-1-4.jpg",
	"3-1-5.jpg",
	"3-1-6.jpg",
	"3-1-7.jpg",
	"3-1-8.jpg",
	"3-1-9.jpg",
	"3-1-10.jpg",
	"3-1-11.jpg",
	"3-1-12.jpg",
	"3-1-13.jpg",
	"3-1-14.jpg",
	"3-1-15.jpg",
	"3-1-16.jpg",

	"5-1-1.jpg",
	"5-1-2.jpg",
	"5-1-3.jpg",
	"5-1-4.jpg",
	"5-1-5.jpg",
	"5-1-6.jpg",
	"5-1-7.jpg",
	"5-1-8.jpg",
	"5-1-9.jpg",
	"5-1-10.jpg",
	"5-1-11.jpg",
	"5-1-12.jpg",
	"5-1-13.jpg",
	"5-1-14.jpg",
	"5-1-15.jpg",
	"5-1-16.jpg",

	"6-1-1.jpg",
	"6-1-2.jpg",
	"6-1-3.jpg",
	"6-1-4.jpg",
	"6-1-5.jpg",
	"6-1-6.jpg",
	"6-1-7.jpg",
	"6-1-8.jpg",
	"6-1-9.jpg",
	"6-1-10.jpg",
	"6-1-11.jpg",
	"6-1-12.jpg",
	"6-1-13.jpg",
	"6-1-14.jpg",
	"6-1-15.jpg",
	"6-1-16.jpg",

	"8-1-1.jpg",
	"8-1-2.jpg",
	"8-1-3.jpg",
	"8-1-4.jpg",
	"8-1-5.jpg",
	"8-1-6.jpg",
	"8-1-7.jpg",
	"8-1-8.jpg",
	"8-1-9.jpg",
	"8-1-10.jpg",
	"8-1-11.jpg",
	"8-1-12.jpg",
	"8-1-13.jpg",
	"8-1-14.jpg",
	"8-1-15.jpg",
	"8-1-16.jpg",

	"9-1-1.jpg",
	"9-1-2.jpg",
	"9-1-3.jpg",
	"9-1-4.jpg",
	"9-1-5.jpg",
	"9-1-6.jpg",
	"9-1-7.jpg",
	"9-1-8.jpg",
	"9-1-9.jpg",
	"9-1-10.jpg",
	"9-1-11.jpg",
	"9-1-12.jpg",
	"9-1-13.jpg",
	"9-1-14.jpg",
	"9-1-15.jpg",
	"9-1-16.jpg",
}

var globalBoxV3 = []boxItemV3{
	{
		1,
		"可爱",
		250,
		670,
		[]string{"伶俐的"},
		[]int{2, 3, 21, 22, 30, 31, 23},
		[]int{7, 20},
		"明黄、橘粉、草绿、浅蓝。色彩明度高、纯度适中，像糖果一样甜美活泼。",
		[]string{"甜美软糯", "治愈亲和", "梦幻童趣"},
		"你的内心住着一个充满童话感的小孩。你对（圆润的线条、马卡龙色系（如蜜桃粉、婴儿蓝、鹅黄）、蓬松的材质）有着天然的偏爱。（生活方式加入）在审美上，你抗拒尖锐和冷漠，更倾向于用甜美、亲和、无害的视觉语言来表达自己。你追求的不仅是外表的萌趣，更是对世界保持一份天真与柔软的向往。",
		"拒绝尖锐，用圆润与甜美拥抱世界。",
	},
	{
		2,
		"闲适",
		130,
		1080,
		[]string{"田园的", "有亲和力的"},
		[]int{1, 10, 11, 12},
		[]int{25, 45},
		"暖黄、草绿、橘色、天蓝。色彩偏暖、自然，像阳光下的草地，让人感觉放松。",
		[]string{"慵懒自在", "亲和无压", "居家度假"},
		"你是一个追求“松弛感”的人。你的审美核心是“舒服自在”，偏爱（棉麻、藤编、原木等自然材质，喜欢宽松、无拘无束的剪裁和带着暖意的色彩）。（生活方式加入）你不喜欢被规则束缚，容易被阳光、微风、下午茶等慵懒的元素打动。在你眼里，最好的风格就是毫无用力感的自然流露。",
		"拒绝内耗，审美就是要绝对的舒服自在。",
	},
	{
		3,
		"动感",
		130,
		1680,
		[]string{"充满活力的"},
		[]int{0},
		[]int{19, 48},
		"正红、明黄、亮蓝、深紫。高饱和度原色与大色块碰撞，视觉节奏快，充满运动活力。",
		[]string{"活力碰撞", "高能冲击", "运动波普"},
		"你是一个充满活力和能量的人。你偏爱高饱和度的对比色（如红绿撞色、亮黄配黑）、波普风格的图案、运动感的剪裁。（生活方式加入）你的审美里有一种不安分的“冲劲”，你喜欢视觉冲击力强、有速度感和节奏感的设计。你觉得色彩和风格就应该用来释放内心的热情与活力，你觉得色彩可以疗愈心灵。",
		"拒绝沉闷，用强烈的色彩碰撞释放内心的能量。",
	},
	{
		4,
		"华丽",
		380,
		1680,
		[]string{"华丽的", "富于装饰的"},
		[]int{91, 92, 100},
		[]int{33, 36},
		"深酒红、暗金、深紫、墨绿。色彩明度极低，饱和度高，带有丝绒和黄金的厚重感。",
		[]string{"繁复丰盛", "复古华丽", "仪式感"},
		"你欣赏那种带有“丰盛感”和“复古底蕴”的美。你偏爱浓郁的色彩（如酒红、墨绿、藏金）、丝绒或金属的质感，以及带有繁复花纹的设计。（生活方式加入）你的审美不属于极简主义，你喜欢被精致且有重量的细节包裹，追求的是一种有仪式感、有历史沉淀的贵气与华丽。",
		"偏爱浓郁与质感，日常也要有华丽的仪式感。",
	},
	{
		5,
		"粗犷",
		380,
		2130,
		[]string{"粗犷的"},
		[]int{101},
		[]int{6},
		"深砖红、铁锈色、黑褐、深咖。色彩极暗、极重，带有泥土和做旧皮革的粗粝感。",
		[]string{"原始力量", "野性不羁", "真实厚重"},
		"你欣赏原始、野性，自然与不加修饰的力量感。你偏爱深棕、暗红、铁锈色，喜欢做旧皮革、粗花呢、岩石等带有耐磨属性的材质。（生活方式加入）你的审美有一种“不受驯服”的野性，你不喜欢过度精细的修饰，反而觉得那些粗糙的纹理、厚重的轮廓，才是真实的、有生命力的表达。",
		"拒绝过度修饰，青睐真实粗粝的原始力量。",
	},

	{
		6,
		"浪漫",
		650,
		620,
		[]string{"浪漫的"},
		[]int{32, 33, 34, 41, 42, 43, 44, 49, 50, 122},
		[]int{10, 27},
		"浅粉、丁香紫、鹅黄、薄荷绿。色彩极度轻盈、通透，带有水彩般的渐变与柔和感。",
		[]string{"轻盈梦幻", "氛围至上", "感性诗性"},
		"你是一个充满想象力和感受力的人。你偏爱轻盈飘逸的材质（如薄纱、蕾丝）、莫兰迪色或梦幻的渐变色。（生活方式加入）你的审美里有一种不切实际的“诗性”，你容易被流线、曲线和充满光影感的事物打动。你不追求实用主义，你更愿意为情绪价值和氛围感买单。",
		"比起实用，你更愿意为浪漫和氛围感买单。",
	},
	{
		7,
		"自然",
		650,
		870,
		[]string{"质朴的"},
		[]int{13, 20, 24, 25, 39, 40, 51, 52, 53, 54, 55, 59, 61, 62, 70, 82, 123},
		[]int{14, 40},
		"各种深浅的棕色、橄榄绿、卡其、泥土色。色彩仿佛直接取自大自然，质朴温润。",
		[]string{"有机大地", "生机平静", "返璞归真"},
		"你拥有很强的“大地情结”。你偏爱植物绿、燕麦色、泥土棕，喜欢一切带有肌理感和原始美感的材质。（生活方式加入）你的审美倾向于生机勃勃又平静的日常，比起人造的精致，你更欣赏那种经历时光洗礼的自然之美。你向往返璞归真的生活状态，追求内心的平静与生机。",
		"与自然共处，让生机自然流淌。",
	},
	{
		8,
		"雅致",
		960,
		1040,
		[]string{"端庄的", "优雅的"},
		[]int{4, 5, 14, 19, 29, 60, 63, 64, 69, 83, 84, 124, 125},
		[]int{12, 15},
		"低饱和度的浅灰粉、灰绿、浅咖与米白。色彩柔和，带有灰蒙蒙的质感，不刺眼。",
		[]string{"柔彩浪漫", "自然闲适", "松弛雅致"},
		"你偏爱一种“松弛的精致感”，但不追求甜腻夸张。你的审美更接近（低饱和的奶油白、雾粉、浅杏、鼠尾草绿与浅木色，喜欢留白、柔光和自然纹理）带来的舒适感。（生活方式加入）你容易被有温度、有呼吸感的美吸引，理想的状态不是刻意装饰，而是浪漫得不费力、自然又雅致地融入日常。",
		"不甜腻，浪漫有呼吸，自然又雅致。",
	},
	{
		9,
		"精致",
		1350,
		1150,
		[]string{"风度翩翩的"},
		[]int{6, 13, 15, 28, 65, 66, 68, 85, 86, 87, 89},
		[]int{1, 29},
		"纯黑、高级灰、藏蓝、暗绿。色彩极简克制，以无彩色为主，充满秩序感。 蒙的质感，不刺眼。",
		[]string{"知性洒脱", "利落高级", "完美主义"},
		"你的理想空间，是一个“藏着高级秩序感”的精神自留地。你偏爱饱和度不高、冷调亮色，你偏爱(颜色）你不喜欢浮夸的装饰堆砌，你追求的精致是“藏”在细节里的：考究的收口、触感极好的材质、精准的灯光氛围》（生活方式加入）你的审美里有着极高的标准，你的精致不是珠光宝气的堆砌，而是克制、精准、每一处细节都经得起推敲的高级感。",
		"在潇洒自在的高智感空间里，做自己。",
	},
	{
		10,
		"古典",
		800,
		1680,
		[]string{"深邃的", "文化气息的"},
		[]int{9, 71, 72, 81, 90, 93, 99, 102, 103},
		[]int{21, 31},
		"做旧的卡其、燕麦色、暗红、褐灰色。色彩像褪色的老照片或复古壁纸，温暖且怀旧。",
		[]string{"复古温度", "经典怀旧", "人文故事"},
		"你更偏向于有“人情味”的复古感。你偏爱复古黄、暗红、做旧的金属，喜欢有年代感的器物和带有复古元素的穿搭。（生活方式加入）你的审美有一种怀旧的温度，你欣赏那些过去岁月里的经典设计，你希望用一种有故事感、有文化底蕴的视觉语言来装点生活。",
		"在复古的余温里，寻找有故事的经典之美。",
	},
	{
		11,
		"考究",
		1150,
		1780,
		[]string{"严谨的"},
		[]int{8, 73, 74, 75, 76, 77, 78, 79, 80, 88, 94, 95, 96, 97, 98, 126, 127},
		[]int{23, 26},
		"深灰蓝、炭灰、暗橄榄绿。色彩偏向冷调暗色，如同建筑图纸或精纺西装，严谨克制。",
		[]string{"理性精准", "秩序结构", "工匠精神"},
		"你更偏向于一种“理性的严谨”。你偏爱冷灰、藏蓝、纯黑，喜欢直线条、对称设计以及极其规整的排版或剪裁。（生活方式加入）你的审美带着一种学术派或工匠式的严谨，你注重材质的功能性与结构的合理性，你欣赏那种一丝不苟、绝对精准、充满秩序感的美学。",
		"严谨有理，迷恋绝对精准的秩序与结构。",
	},
	{
		12,
		"格调",
		860,
		2130,
		[]string{"有格调的", "独特的"},
		[]int{104, 105, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 128, 129, 130},
		[]int{8, 38},
		"暗红、墨绿、深棕、藏青。色彩仿佛老旧图书馆里的皮质沙发，沉稳且充满底蕴。",
		[]string{"历史沉淀", "庄重严谨", "人文厚度"},
		"你是一个敬畏时间与历史的人。喜欢独特的设计，你偏爱深色调（如勃艮第红、藏青、深橄榄绿），喜欢有结构感、有厚度、的设计。（生活方式加入）你的审美融合了古典的庄重与考究的严谨，你觉得经得起时间考验的美一定是充满沉淀的、有规矩的、带着人文厚度的。",
		"敬畏时间，欣赏充满规矩与沉淀的厚重之美。",
	},
	{
		13,
		"正式",
		1440,
		2100,
		[]string{"正统的"},
		[]int{106, 107},
		[]int{47},
		"深藏青、炭灰、黑色。色彩明度极低且纯度低，是最经典的商务西装色，权威稳重。",
		[]string{"权威力量", "专业得体", "气场强大"},
		"你偏好极具“权威感”和“专业感”的风格。你偏爱深黑、深灰、藏青等稳重色，喜欢极具结构感的设计。（生活方式加入）你的审美里有一种强大的气场，你觉得空间和色彩是用来展现力量、专业度和边界感的工具，你追求的是无懈可击的得体与气场。",
		"用稳重的轮廓和低调的色彩，构筑强大的气场。",
	},
	{
		14,
		"清爽",
		1590,
		730,
		[]string{"清净的"},
		[]int{26, 35, 36, 37, 38, 45, 46, 47, 48, 56, 57, 58, 121},
		[]int{28, 37},
		"水蓝色、薄荷绿、纯白。色彩明度高、冷感强，仿佛夏日微风或冰镇气泡水。",
		[]string{"干净通透", "明亮轻盈", "极简纯净"},
		"你偏爱那种“如沐春风”的干净感。你的审美极其纯净，偏爱蓝白调、透明材质、玻璃感以及极简的线条。你讨厌繁复和压抑，喜欢呼吸感强的视觉体验。（生活方式加入）你的审美就像夏日的微风或刚洗过的白衬衫，给人带来一种明亮、通透、毫无杂质的轻盈感。",
		"如微风拂面，用干净和通透治愈一切。",
	},
	{
		15,
		"青春",
		1650,
		1160,
		[]string{"青春洋溢的"},
		[]int{16, 17, 27, 67},
		[]int{17, 41},
		"灰蓝、冷绿、深水蓝，点缀小面积亮黄。色彩偏冷、灰度高，带有北欧的疏离与清净。",
		[]string{"清冷疏离", "北欧极简", "独立松弛"},
		"你追求的是一种“有距离感的松弛”。你偏爱冷灰色、雾霾蓝、冰川白，喜欢北欧极简风、利落的线条和克制的材质。你不喜欢喧闹，渴望在纷扰的世界中保持一份清冷和独立。（生活方式加入）你的闲适，是即使穿着最简单的冷淡风单品，也能散发出一种不费力的疏离感。",
		"保持一点距离感，在清冷中享受不费力的松弛。",
	},
	{
		16,
		"现代",
		1690,
		1690,
		[]string{"都市气息的", "革新的"},
		[]int{7},
		[]int{4, 35},
		"大面积黑、白、深灰，点缀极高饱和度的亮绿、明黄。视觉冲击力强，有科技感。",
		[]string{"前卫几何", "都市利落", "科技先锋"},
		"你有一种敏锐的“都市感”。你偏爱无彩色系（黑、白、灰）和冷色调，喜欢搭配极具视觉冲击力的高饱和亮色，喜欢几何图形、不对称设计和科技感材质。（生活方式加入）你的审美非常前瞻，你不喜欢怀旧，更愿意拥抱变化。你追求的是利落、高效、充满未来感和不设限的视觉风格。",
		"告别陈旧，用利落与几何拥抱未来感。",
	},
}

func (c *Controller) GetColorbarListV3(ctx *gin.Context) {
	colorbars := c.service.GetIndexColorbar()

	c.ResponseSuccess(ctx, colorbars)
}

func (c *Controller) GetImageListV3(ctx *gin.Context) {
	gender := ctx.Query("gender")
	imageSource := globalImageV3

	imageData := make(map[string][]map[string]interface{})
	for _, category := range globalImageCategory {
		images := make([]map[string]interface{}, 0)
		for _, imageName := range imageSource {
			parts := strings.Split(imageName, "-")
			if len(parts) < 3 {
				continue
			}
			categoryID := cast.ToInt(parts[0])
			if categoryID != category.CategoryID {
				continue
			}
			genderPart := parts[1]
			if categoryID == 1 && gender == "男" && genderPart != "2" {
				continue
			}
			if categoryID == 1 && gender == "女" && genderPart != "1" {
				continue
			}

			images = append(images, map[string]interface{}{
				"id":           imageName,
				"url":          kcfg.GetString("app.global.host") + "/imgv3/" + imageName,
				"categoryName": imageCategoryIDToName(categoryID),
			})
		}
		imageData[category.CategoryEnglishName] = images
	}

	c.ResponseSuccess(ctx, imageData)
}

func (c *Controller) GetWordListV3(ctx *gin.Context) {
	words := getAllWordV3()
	c.ResponseSuccess(ctx, words)
}

type ColorbarDetailV3 struct {
	Num int    `json:"num"`
	URL string `json:"url"`
}

type AestheticDataRspV3 struct {
	AestheticDataRsp
	Tag                []string           `json:"tag"`
	Attitude           string             `json:"attitude"`
	ColorBars          []ColorbarDetailV3 `json:"colorbars"`
	Background         string             `json:"background"`
	LikedColorWords    []string           `json:"liked_color_words"`
	DislikedColorWords []string           `json:"disliked_color_words"`
}

func (c *Controller) GetUserAestheticDataListV3(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	resp, err := c.service.GetUserAestheticDataListV3(c.GetUserID(ctx), page, pageSize)
	if err != nil {
		c.ResponseError(ctx, 500, "获取数据失败: "+err.Error())
		return
	}

	c.ResponseSuccess(ctx, resp)
}

func (c *Controller) GetAestheticDataDetailV3(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.ResponseError(ctx, 400, "无效的ID参数")
		return
	}
	data, err := c.service.GetAestheticDataDetailV3(uint(id), c.GetUserID(ctx), false)
	if err != nil {
		c.ResponseError(ctx, 500, "获取数据失败: "+err.Error())
		return
	}
	c.ResponseSuccess(ctx, data)
}

// v3列表获取全部数据，而不是只获取version为3的数据
func (s *Service) GetUserAestheticDataListV3(userID uint, page, pageSize int) (*PageResponse, error) {
	var total int64
	var list []AestheticData

	if err := s.db.Model(&AestheticData{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, err
	}
	if err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error; err != nil {
		return nil, err
	}

	result := make([]*AestheticDataRsp, 0, len(list))
	for k, v := range list {
		list[k].ColorImageURL = getColorImageURL(v.ColorImageURL)
		list[k].BoxImageURL = getBoxImageURL(v.BoxImageURL)

		likedColors, likedColorBars := []int{}, []int{}
		likedAdjectives, likedImages := []string{}, []string{}
		_ = json.Unmarshal([]byte(v.LikedColors), &likedColors)
		_ = json.Unmarshal([]byte(v.LikedColorbars), &likedColorBars)
		_ = json.Unmarshal([]byte(v.LikedAdjectives), &likedAdjectives)
		_ = json.Unmarshal([]byte(v.LikedImages), &likedImages)

		_, primaryBox, _, _ := buildAestheticSummaryV3(likedColors, likedColorBars, likedAdjectives, likedImages)
		comment := make([]string, 0, 1)
		if primaryBox.Comment != "" {
			comment = append(comment, primaryBox.Comment)
		}

		result = append(result, &AestheticDataRsp{
			AestheticData: list[k],
			Comment:       comment,
		})
	}

	return &PageResponse{List: result, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) GetAestheticDataDetailV3(id, userID uint, showSolutionDetal bool) (*AestheticDataRspV3, error) {
	var data AestheticData
	if err := s.db.Where("id = ? AND version = ?", id, 3).First(&data).Error; err != nil {
		return nil, err
	}
	data.ColorImageURL = getColorImageURL(data.ColorImageURL)
	data.BoxImageURL = getBoxImageURL(data.BoxImageURL)
	likedColors, dislikedColors, likedColorBars := []int{}, []int{}, []int{}
	likedAdjectives, likedImages := []string{}, []string{}
	_ = json.Unmarshal([]byte(data.LikedColors), &likedColors)
	_ = json.Unmarshal([]byte(data.DislikedColors), &dislikedColors)
	_ = json.Unmarshal([]byte(data.LikedColorbars), &likedColorBars)
	_ = json.Unmarshal([]byte(data.LikedAdjectives), &likedAdjectives)
	_ = json.Unmarshal([]byte(data.LikedImages), &likedImages)
	likedLifeImageDesc, likedAreaImageDesc := make([]string, 0, len(likedImages)), make([]string, 0)
	for idx, imageName := range likedImages {
		likedLifeImageDesc = append(likedLifeImageDesc, getGlobalImageURLV3(imageName))
		if idx >= 4 {
			likedAreaImageDesc = append(likedAreaImageDesc, getGlobalImageURLV3(imageName))
		}
	}
	summary, primaryBox, secondaryBox, resultWord := buildAestheticSummaryV3(likedColors, likedColorBars, likedAdjectives, likedImages)
	comment := make([]string, 0, 1)
	if primaryBox.Comment != "" {
		comment = append(comment, primaryBox.Comment)
	}
	solutionList := make([]SolutionItem, 0)
	if showSolutionDetal && primaryBox.Num > 0 {
		solutionList = append(solutionList, SolutionItem{Icon: getSolutionIconV3(primaryBox.Num), Images: getSolutionImagesV3(primaryBox.Num)})
		if secondaryBox.Num > 0 && secondaryBox.Num != primaryBox.Num {
			solutionList = append(solutionList, SolutionItem{Icon: getSolutionIconV3(secondaryBox.Num), Images: getSolutionImagesV3(secondaryBox.Num)})
		}
	}
	isCollection, targetSolutionID, collectionTotal := 0, 0, int64(0)
	if err := s.db.Model(&UserCollection{}).Where("user_id = ?", userID).Where("a_id = ?", cast.ToInt(data.ID)).Count(&collectionTotal).Error; err == nil && collectionTotal >= 1 {
		isCollection = 1
	}
	var targetSolution UserSolution
	if err := s.db.Model(&UserSolution{}).Where("user_id = ?", userID).Where("a_id = ?", cast.ToInt(data.ID)).First(&targetSolution).Error; err == nil {
		targetSolutionID = int(targetSolution.ID)
	}
	background := getBackgroundURLV3(primaryBox.Num)
	likedColorWords := buildTopColorWordsV3(likedColors, 6)
	dislikedColorWords := buildTopColorWordsV3(dislikedColors, 6)
	return &AestheticDataRspV3{AestheticDataRsp: AestheticDataRsp{AestheticData: data, Comment: comment, LikedColorDesc: buildColorDescItemsV3(likedColors), DislikedColorDesc: buildColorDescItemsV3(dislikedColors), LikedAdjectiveDesc: likedAdjectives, LikedLifeImageDesc: likedLifeImageDesc, LikedAreaImageDesc: likedAreaImageDesc, Summary: summary, WordBaseOnColor: resultWord, IsCollection: isCollection, SolutionList: solutionList, TargetSolutionID: targetSolutionID}, Tag: primaryBox.Tag, Attitude: primaryBox.Attitude, ColorBars: getColorBarDetailsV3(primaryBox.ColorBars), Background: background, LikedColorWords: likedColorWords, DislikedColorWords: dislikedColorWords}, nil
}

func getBackgroundURLV3(boxNum int) string {
	return kcfg.GetString("app.global.host") + "/background/" + cast.ToString(boxNum) + ".jpg"
}

type IndexColorBarRspItem struct {
	Num int    `json:"num"`
	URL string `json:"url"`
}

func (s *Service) GetIndexColorbar() []*IndexColorBarRspItem {
	rsp := make([]*IndexColorBarRspItem, 0)
	for _, v := range globalColorbar {
		rsp = append(rsp, &IndexColorBarRspItem{
			Num: v.Num,
			URL: kcfg.GetString("app.global.host") + v.URL,
		})
	}
	return rsp
}

type AestheticDataRequestV3 struct {
	Name            string   `json:"name" binding:"required"`                         // 姓名
	Gender          string   `json:"gender" binding:"required"`                       // 性别
	Age             int      `json:"age" binding:"required,min=1,max=120"`            // 年龄
	City            string   `json:"city" binding:"required"`                         // 所在城市
	Phone           string   `json:"phone" binding:"required"`                        // 手机号码
	LikedColors     []string `json:"liked_colors" binding:"required,min=1,max=10"`    // 喜欢的颜色
	DislikedColors  []string `json:"disliked_colors" binding:"required,min=1,max=5"`  // 讨厌的颜色
	LikedAdjectives []string `json:"liked_adjectives" binding:"required,min=1,max=3"` // 喜欢的形容词
	LikedColorBars  []string `json:"liked_colorbars" binding:"required,min=1,max=3"`  // 喜欢的颜色条
	LikedImages     []string `json:"liked_images" binding:"required,min=1,max=12"`    // 喜欢的图片
}

func (c *Controller) SaveAestheticDataV3(ctx *gin.Context) {
	var req AestheticDataRequestV3
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.ResponseError(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	// 本地串行化高开销的保存流程，避免瞬时高并发压垮当前服务进程。
	select {
	case c.saveAestheticDataLimiter <- struct{}{}:
		defer func() {
			<-c.saveAestheticDataLimiter
		}()
	case <-ctx.Request.Context().Done():
		c.ResponseError(ctx, 499, "请求已取消，请稍后重试")
		return
	}

	userID := c.GetUserID(ctx)
	data, err := c.service.SaveAestheticDataV3(userID, c.GetUserSource(ctx), &req)
	if err != nil {
		c.ResponseError(ctx, 500, "保存数据失败: "+err.Error())
		return
	}

	// 返回保存的数据
	c.ResponseSuccess(ctx, data)
}

func (s *Service) SaveAestheticDataV3(userID uint, source int, req *AestheticDataRequestV3) (*AestheticDataRsp, error) {
	// 获取用户信息
	var user User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	likedColorNum := make([]int, 0)
	for _, v := range req.LikedColors {
		// 如果没有#说明是色号，不需要转换
		if strings.Index(v, "#") == -1 {
			likedColorNum = append(likedColorNum, cast.ToInt(v))
		} else {
			// 如果有#说明是rgb值，需要转换
			likedColorNum = append(likedColorNum, rgbToNum(v))
		}
	}

	disLikedColorNum := make([]int, 0)
	for _, v := range req.DislikedColors {
		if strings.Index(v, "#") == -1 {
			disLikedColorNum = append(disLikedColorNum, cast.ToInt(v))
		} else {
			disLikedColorNum = append(disLikedColorNum, rgbToNum(v))
		}
	}
	likedColorbarsNum := make([]int, 0)
	for _, v := range req.LikedColorBars {
		likedColorbarsNum = append(likedColorbarsNum, cast.ToInt(v))
	}

	// 序列化数组为JSON字符串
	likedColorsJSON, err := json.Marshal(likedColorNum)
	if err != nil {
		return nil, err
	}

	dislikedColorsJSON, err := json.Marshal(disLikedColorNum)
	if err != nil {
		return nil, err
	}

	likedColorbarsJSON, err := json.Marshal(likedColorbarsNum)
	if err != nil {
		return nil, err
	}

	likedAdjectivesJSON, err := json.Marshal(req.LikedAdjectives)
	if err != nil {
		return nil, err
	}

	likedImagesJSON, err := json.Marshal(req.LikedImages)
	if err != nil {
		return nil, err
	}

	// 绘制颜色图片
	colorimg := DrawColorV3(likedColorNum, disLikedColorNum)
	// 绘制坐标图片
	boximg := DrawToBoxV3(req.LikedImages, likedColorNum, disLikedColorNum, req.LikedAdjectives)

	// 创建审美数据记录
	data := AestheticData{
		UserID:          userID,
		Name:            req.Name,
		Gender:          req.Gender,
		Age:             req.Age,
		City:            req.City,
		Phone:           req.Phone,
		LikedColors:     string(likedColorsJSON),
		DislikedColors:  string(dislikedColorsJSON),
		LikedColorbars:  string(likedColorbarsJSON),
		LikedAdjectives: string(likedAdjectivesJSON),
		LikedImages:     string(likedImagesJSON),
		ColorImageURL:   trimBasePath(colorimg),
		BoxImageURL:     trimBasePath(boximg),
		Source:          source,
		Version:         3, // 新版本使用3
	}

	// 保存数据
	if err := s.db.Create(&data).Error; err != nil {
		return nil, err
	}

	_, primaryBox, _, _ := buildAestheticSummaryV3(likedColorNum, likedColorbarsNum, req.LikedAdjectives, req.LikedImages)
	comment := make([]string, 0, 1)
	if primaryBox.Comment != "" {
		comment = append(comment, primaryBox.Comment)
	}
	// 返回保存的数据
	return &AestheticDataRsp{
		AestheticData:    data,
		ColorExplorImage: getColorImageURL(colorimg),
		BoxExplorImage:   getBoxImageURL(boximg),
		Comment:          comment,
	}, nil
}

var (
	renderIndexOnceV3  sync.Once
	boxByNumV3         map[int]boxItemV3
	colorToBoxNumV3    map[int]int
	wordToBoxNumV3     map[string]int
	colorbarToBoxNumV3 map[int]int
)

func ensureRenderIndexesV3() {
	renderIndexOnceV3.Do(func() {
		boxByNumV3 = make(map[int]boxItemV3, len(globalBoxV3))
		colorToBoxNumV3 = make(map[int]int)
		wordToBoxNumV3 = make(map[string]int)
		colorbarToBoxNumV3 = make(map[int]int)

		for _, box := range globalBoxV3 {
			boxByNumV3[box.Num] = box
			for _, colorNum := range box.Colors {
				colorToBoxNumV3[colorNum] = box.Num
			}
			for _, word := range box.Words {
				wordToBoxNumV3[word] = box.Num
			}
			for _, colorbarNum := range box.ColorBars {
				colorbarToBoxNumV3[colorbarNum] = box.Num
			}
		}
	})
}

func DrawToBoxV3(likedImages []string, likedColor, dislikedColor []int, words []string) string {
	ensureRenderIndexesV3()

	tempDir := "./"

	// 输入图片路径 (实际使用时请替换为实际图片路径)
	inputPath := filepath.Join(tempDir, "boxbase-5.jpg")

	// 输出图片路径
	outputPath := filepath.Join(tempDir, fmt.Sprintf("boximg/outputV3_%d_%d.jpg", time.Now().UnixNano(), rand.Intn(1000)+9000))

	backgroundImage, err := loadImageAsRGBA(inputPath)
	if err != nil {
		return ""
	}

	// 1.把图片合并到抽屉图

	// 3.把喜欢的颜色合并到抽屉图
	likedcolormap := make(map[int][]int, 0)
	for _, v := range likedColor {
		if boxNum, ok := colorToBoxNumV3[v]; ok {
			likedcolormap[boxNum] = append(likedcolormap[boxNum], v)
		}
	}
	for boxNum, v := range likedcolormap {
		box, ok := boxByNumV3[boxNum]
		if !ok {
			continue
		}

		startX := box.StartX
		startY := box.StartY
		for _, col := range v {
			// 画喜欢的颜色，x不变，y增加250
			curY := startY + 10
			r, g, b := numToRGB(col)
			drawSolidRectangleRGBA(backgroundImage, startX, curY, 100, 100, color.RGBA{R: r, G: g, B: b, A: 255})
			// 如果是多个颜色，x增加100
			startX = startX + 100
		}
	}

	// 4.把不喜欢的颜色合并到抽屉图
	dislikedcolormap := make(map[int][]int, 0)
	for _, v := range dislikedColor {
		if boxNum, ok := colorToBoxNumV3[v]; ok {
			dislikedcolormap[boxNum] = append(dislikedcolormap[boxNum], v)
		}
	}
	for boxNum, v := range dislikedcolormap {
		box, ok := boxByNumV3[boxNum]
		if !ok {
			continue
		}

		startX := box.StartX
		startY := box.StartY
		for _, col := range v {
			// 画不喜欢的颜色，x不变，y增加250
			curY := startY + 110
			r, g, b := numToRGB(col)
			drawCrossRGBA(backgroundImage, startX, curY, 100, 20, color.RGBA{R: r, G: g, B: b, A: 255})
			// 如果是多个颜色，x增加55
			startX = startX + 100
		}
	}

	// 把关键词合并到抽屉图
	wordmap := make(map[int][]string, 0)
	for _, v := range words {
		if boxNum, ok := wordToBoxNumV3[v]; ok {
			wordmap[boxNum] = append(wordmap[boxNum], v)
		}
	}
	for vnum, v := range wordmap {
		box, ok := boxByNumV3[vnum]
		if !ok {
			continue
		}

		startX := box.StartX
		startY := box.StartY
		curY := startY + 230
		tmpstring := strings.Join(v, ",")
		if err := drawLabelTextRGBA(backgroundImage, startX, curY, 60, tmpstring, color.RGBA{R: 86, G: 102, B: 169, A: 255}, "./msyhbd.ttc"); err != nil {
			return ""
		}
		// 使用freetype将文字写入图片overlap中
	}

	if err := saveRenderedImage(outputPath, backgroundImage); err != nil {
		return ""
	}
	return outputPath
}

func DrawColorV3(likedColorNum, disLikedColorNum []int) string {
	// 获取临时目录
	tempDir := "./"

	// 输入图片路径 (实际使用时请替换为实际图片路径)
	inputPath := filepath.Join(tempDir, "colorbase-4.jpg")

	// 输出图片路径
	outputPath := filepath.Join(tempDir, fmt.Sprintf("colorimg/outputV3_%d_%d.jpg", time.Now().UnixNano(), rand.Intn(1000)+9000))

	backgroundImage, err := loadImageAsRGBA(inputPath)
	if err != nil {
		return ""
	}

	for _, v := range likedColorNum {
		xstep := ((v - 1) % 10)
		ystep := (v - 1) / 10
		// 设置矩形参数
		x := 210 + xstep*87 // 矩形左上角x坐标
		y := 306 + ystep*88 // 矩形左上角y坐标
		width := 88         // 矩形宽度
		height := 88        // 矩形高度
		r, g, b := numToRGB(v)
		if ystep >= 12 {
			y = y + 15
		} else if ystep > 8 {
			y = y + 5
		} else if ystep > 6 {
			y = y + 2
		}

		// 画喜欢的颜色，使用矩形
		drawSolidRectangleRGBA(backgroundImage, x, y, width, height, color.RGBA{R: r, G: g, B: b, A: 255})
	}

	for _, v := range disLikedColorNum {

		xstep := ((v - 1) % 10)
		ystep := (v - 1) / 10
		// 设置矩形参数
		x := 210 + xstep*87 // 矩形左上角x坐标
		y := 306 + ystep*88 // 矩形左上角y坐标
		width := 88         // 三角形宽度
		r, g, b := numToRGB(v)
		if ystep >= 12 {
			y = y + 15
		} else if ystep > 8 {
			y = y + 5
		} else if ystep > 6 {
			y = y + 2
		}

		// 画讨厌的颜色，使用X
		drawCrossRGBA(backgroundImage, x, y, width, 20, color.RGBA{R: r, G: g, B: b, A: 255})
	}

	if err := saveRenderedImage(outputPath, backgroundImage); err != nil {
		fmt.Printf("处理图片时出错: %v\n", err)
		return ""
	}

	return outputPath
}
