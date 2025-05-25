// 生成随机数据
function generateMockData(count) {
    const regions = ['北京', '上海', '广州', '深圳', '杭州', '成都', '武汉', '西安', '南京', '重庆'];
    const colors = ['红色', '蓝色', '绿色', '黄色', '紫色', '粉色', '白色', '黑色', '橙色', '灰色'];
    const adjectives = ['简约', '现代', '复古', '奢华', '文艺', '可爱', '优雅', '时尚', '清新', '典雅'];
    const surnames = ['张', '王', '李', '赵', '刘', '陈', '杨', '黄', '周', '吴'];
    
    return Array.from({ length: count }, (_, index) => {
        // 随机生成2-4个喜欢的颜色
        const likedColorsCount = Math.floor(Math.random() * 3) + 2;
        const likedColors = shuffle([...colors]).slice(0, likedColorsCount);
        
        // 随机生成1-3个讨厌的颜色（不包含喜欢的颜色）
        const availableDislikeColors = colors.filter(color => !likedColors.includes(color));
        const dislikedColorsCount = Math.floor(Math.random() * 3) + 1;
        const dislikedColors = shuffle(availableDislikeColors).slice(0, dislikedColorsCount);
        
        // 随机生成1-3个喜欢的形容词
        const likedAdjectivesCount = Math.floor(Math.random() * 3) + 1;
        const likedAdjectives = shuffle([...adjectives]).slice(0, likedAdjectivesCount);

        return {
            name: surnames[Math.floor(Math.random() * surnames.length)] + 
                  ['小明', '小红', '小华', '小芳', '小强', '小美', '小龙', '小云'][Math.floor(Math.random() * 8)],
            age: Math.floor(Math.random() * 42) + 18, // 18-60岁
            gender: Math.random() > 0.5 ? 'male' : 'female',
            region: regions[Math.floor(Math.random() * regions.length)],
            likedColors,
            dislikedColors,
            likedAdjectives
        };
    });
}

// Fisher-Yates 洗牌算法
function shuffle(array) {
    for (let i = array.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [array[i], array[j]] = [array[j], array[i]];
    }
    return array;
}

// 生成100条模拟数据
const mockData = generateMockData(100);

// 全局变量
let currentPage = 1;
const pageSize = 10;
let filteredData = [...mockData];

// 初始化
document.addEventListener('DOMContentLoaded', () => {
    initRegionOptions();
    renderTable();
    renderCharts();
    initTabs();
});

// 初始化地区选项
function initRegionOptions() {
    const regions = [...new Set(mockData.map(item => item.region))];
    const regionSelect = document.getElementById('region');
    regions.forEach(region => {
        const option = document.createElement('option');
        option.value = region;
        option.textContent = region;
        regionSelect.appendChild(option);
    });
}

// 应用筛选
function applyFilters() {
    const ageRange = document.getElementById('ageRange').value;
    const gender = document.getElementById('gender').value;
    const region = document.getElementById('region').value;

    filteredData = mockData.filter(item => {
        const ageMatch = !ageRange || checkAgeRange(item.age, ageRange);
        const genderMatch = !gender || item.gender === gender;
        const regionMatch = !region || item.region === region;
        return ageMatch && genderMatch && regionMatch;
    });

    currentPage = 1;
    renderTable();
    renderCharts();
}

// 检查年龄范围
function checkAgeRange(age, range) {
    switch(range) {
        case '0-18': return age <= 18;
        case '19-30': return age > 18 && age <= 30;
        case '31-50': return age > 30 && age <= 50;
        case '50+': return age > 50;
        default: return true;
    }
}

// 渲染表格
function renderTable() {
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = startIndex + pageSize;
    const pageData = filteredData.slice(startIndex, endIndex);

    const tbody = document.getElementById('dataTableBody');
    tbody.innerHTML = '';

    pageData.forEach(item => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${item.name}</td>
            <td>${item.age}</td>
            <td>${item.gender === 'male' ? '男' : '女'}</td>
            <td>${item.region}</td>
            <td>${item.likedColors.join(', ')}</td>
            <td>${item.dislikedColors.join(', ')}</td>
            <td>${item.likedAdjectives.join(', ')}</td>
        `;
        tbody.appendChild(row);
    });

    renderPagination();
}

// 渲染分页
function renderPagination() {
    const totalPages = Math.ceil(filteredData.length / pageSize);
    const pagination = document.getElementById('pagination');
    pagination.innerHTML = '';

    for (let i = 1; i <= totalPages; i++) {
        const button = document.createElement('button');
        button.textContent = i;
        button.onclick = () => {
            currentPage = i;
            renderTable();
        };
        if (i === currentPage) {
            button.style.backgroundColor = '#007bff';
            button.style.color = 'white';
        }
        pagination.appendChild(button);
    }
}

// 渲染图表
function renderCharts() {
    // 只渲染当前激活的tab的图表
    const activeTab = document.querySelector('.tab-button.active').getAttribute('data-tab');
    if (activeTab === 'color') {
        renderColorChart('like');
        renderColorChart('dislike');
        renderColorPie('like');
        renderColorPie('dislike');
    } else if (activeTab === 'adjective') {
        renderAdjectiveChart();
        renderAdjectivePie();
    } else if (activeTab === 'image') {
        renderBarChart('imageChart', imageData, '喜欢的图片 Top 10');
        renderPieChart('imagePie', imagePercentData, '图片偏好比例');
    } else if (activeTab === 'region') {
        renderRegionMap();
        renderCityChart();
    }
}

// 渲染颜色统计折线图
function renderColorChart(type) {
    const chartDom = document.getElementById(`color${type.charAt(0).toUpperCase() + type.slice(1)}Chart`);
    const myChart = echarts.init(chartDom);

    const colorData = getColorStatistics(type);
    const option = {
        title: {
            text: type === 'like' ? '最喜欢的颜色TOP5' : '最讨厌的颜色TOP5'
        },
        xAxis: {
            type: 'category',
            data: colorData.map(item => item.color)
        },
        yAxis: {
            type: 'value'
        },
        series: [{
            data: colorData.map(item => item.count),
            type: 'line'
        }]
    };

    myChart.setOption(option);
}

// 渲染颜色统计饼图
function renderColorPie(type) {
    const chartDom = document.getElementById(`color${type.charAt(0).toUpperCase() + type.slice(1)}Pie`);
    const myChart = echarts.init(chartDom);

    const colorData = getColorStatistics(type);
    const option = {
        title: {
            text: type === 'like' ? '喜欢颜色占比' : '讨厌颜色占比',
            left: 'center',
            top: 0
        },
        tooltip: {
            trigger: 'item',
            formatter: '{a} <br/>{b}: {c} ({d}%)'
        },
        legend: {
            orient: 'vertical',
            left: 'left',
            top: 'middle'
        },
        grid: {
            top: '10%'
        },
        series: [{
            name: type === 'like' ? '喜欢颜色' : '讨厌颜色',
            type: 'pie',
            radius: ['40%', '70%'],
            center: ['55%', '60%'],
            avoidLabelOverlap: false,
            label: {
                show: true,
                formatter: '{b}\n{d}%',
                position: 'outside'
            },
            emphasis: {
                label: {
                    show: true,
                    fontSize: '16',
                    fontWeight: 'bold'
                }
            },
            labelLine: {
                show: true
            },
            data: colorData.map(item => ({
                name: item.color,
                value: item.count
            }))
        }]
    };

    myChart.setOption(option);
}

// 获取颜色统计数据
function getColorStatistics(type) {
    const colorMap = new Map();
    
    filteredData.forEach(item => {
        const colors = type === 'like' ? item.likedColors : item.dislikedColors;
        colors.forEach(color => {
            colorMap.set(color, (colorMap.get(color) || 0) + 1);
        });
    });

    return Array.from(colorMap.entries())
        .map(([color, count]) => ({ color, count }))
        .sort((a, b) => b.count - a.count)
        .slice(0, 5);
}

// 添加形容词统计折线图
function renderAdjectiveChart() {
    const chartDom = document.getElementById('adjectiveChart');
    const myChart = echarts.init(chartDom);

    const adjectiveData = getAdjectiveStatistics();
    const option = {
        title: {
            text: '最受欢迎形容词TOP5'
        },
        tooltip: {
            trigger: 'axis'
        },
        xAxis: {
            type: 'category',
            data: adjectiveData.map(item => item.adjective)
        },
        yAxis: {
            type: 'value',
            name: '选择人数'
        },
        series: [{
            data: adjectiveData.map(item => item.count),
            type: 'line',
            smooth: true,
            markPoint: {
                data: [
                    {type: 'max', name: '最大值'},
                    {type: 'min', name: '最小值'}
                ]
            }
        }]
    };

    myChart.setOption(option);
}

// 添加形容词统计饼图
function renderAdjectivePie() {
    const chartDom = document.getElementById('adjectivePie');
    const myChart = echarts.init(chartDom);

    const adjectiveData = getAdjectiveStatistics();
    const option = {
        title: {
            text: '形容词偏好占比',
            left: 'center',
            top: 0
        },
        tooltip: {
            trigger: 'item',
            formatter: '{a} <br/>{b}: {c} ({d}%)'
        },
        legend: {
            orient: 'vertical',
            left: 'left',
            top: 'middle'
        },
        grid: {
            top: '10%'
        },
        series: [{
            name: '形容词偏好',
            type: 'pie',
            radius: ['40%', '70%'],
            center: ['55%', '60%'],
            avoidLabelOverlap: false,
            label: {
                show: true,
                formatter: '{b}\n{d}%',
                position: 'outside'
            },
            emphasis: {
                label: {
                    show: true,
                    fontSize: '16',
                    fontWeight: 'bold'
                }
            },
            labelLine: {
                show: true
            },
            data: adjectiveData.map(item => ({
                name: item.adjective,
                value: item.count
            }))
        }]
    };

    myChart.setOption(option);
}

// 获取形容词统计数据
function getAdjectiveStatistics() {
    const adjectiveMap = new Map();
    
    filteredData.forEach(item => {
        item.likedAdjectives.forEach(adj => {
            adjectiveMap.set(adj, (adjectiveMap.get(adj) || 0) + 1);
        });
    });

    return Array.from(adjectiveMap.entries())
        .map(([adjective, count]) => ({ adjective, count }))
        .sort((a, b) => b.count - a.count)
        .slice(0, 5);
}

// 添加地域分析地图
function renderRegionMap() {
    const chartDom = document.getElementById('regionMap');
    const myChart = echarts.init(chartDom);

    // 获取地域统计数据
    const regionData = getRegionStatistics();

    const option = {
        title: {
            text: '用户地域分布',
            left: 'center'
        },
        tooltip: {
            trigger: 'item',
            formatter: '{b}: {c}人'
        },
        visualMap: {
            min: 0,
            max: Math.max(...regionData.map(item => item.value)),
            left: 'left',
            top: 'bottom',
            text: ['高', '低'],
            calculable: true,
            inRange: {
                color: ['#e0f3f8', '#045a8d']
            }
        },
        series: [{
            name: '用户数量',
            type: 'map',
            map: 'china',
            emphasis: {
                label: {
                    show: true
                }
            },
            data: regionData
        }]
    };

    // 加载中国地图数据
    fetch('/admin/areas_full')
        .then(response => response.json())
        .then(geoJson => {
            echarts.registerMap('china', geoJson);
            myChart.setOption(option);
        });
}

// 获取地域统计数据
function getRegionStatistics() {
    const regionMap = new Map();
    const regionNameMap = {
        '北京': '北京市',
        '上海': '上海市',
        '广州': '广东省',
        '深圳': '广东省',
        '杭州': '浙江省',
        '成都': '四川省',
        '武汉': '湖北省',
        '西安': '陕西省',
        '南京': '江苏省',
        '重庆': '重庆市'
    };
    
    filteredData.forEach(item => {
        const provinceName = regionNameMap[item.region];
        regionMap.set(provinceName, (regionMap.get(provinceName) || 0) + 1);
    });

    return Array.from(regionMap.entries()).map(([name, value]) => ({
        name,
        value
    }));
}

// 添加城市统计柱状图
function renderCityChart() {
    const chartDom = document.getElementById('cityChart');
    const myChart = echarts.init(chartDom);

    const cityData = getCityStatistics();
    
    const option = {
        title: {
            text: '用户数量TOP5城市',
            left: 'center'
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'shadow'
            }
        },
        grid: {
            left: '3%',
            right: '4%',
            bottom: '3%',
            containLabel: true
        },
        xAxis: {
            type: 'category',
            data: cityData.map(item => item.city),
            axisLabel: {
                interval: 0,
                rotate: 30
            }
        },
        yAxis: {
            type: 'value',
            name: '用户数量'
        },
        series: [{
            name: '用户数量',
            type: 'bar',
            data: cityData.map(item => item.count),
            itemStyle: {
                color: '#3498db'
            },
            label: {
                show: true,
                position: 'top'
            }
        }]
    };

    myChart.setOption(option);
}

// 获取城市统计数据
function getCityStatistics() {
    const cityMap = new Map();
    
    filteredData.forEach(item => {
        cityMap.set(item.region, (cityMap.get(item.region) || 0) + 1);
    });

    return Array.from(cityMap.entries())
        .map(([city, count]) => ({ city, count }))
        .sort((a, b) => b.count - a.count)
        .slice(0, 5);
}

// 修改initTabs函数
function initTabs() {
    const tabButtons = document.querySelectorAll('.tab-button');
    const tabContents = document.querySelectorAll('.tab-content');

    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            // 移除所有active类
            tabButtons.forEach(btn => btn.classList.remove('active'));
            tabContents.forEach(content => content.classList.remove('active'));

            // 添加active类到当前选中的tab
            button.classList.add('active');
            const tabId = button.getAttribute('data-tab') + 'Tab';
            document.getElementById(tabId).classList.add('active');

            // 重新渲染当前tab的图表
            if (button.getAttribute('data-tab') === 'color') {
                renderColorChart('like');
                renderColorChart('dislike');
                renderColorPie('like');
                renderColorPie('dislike');
            } else if (button.getAttribute('data-tab') === 'adjective') {
                renderAdjectiveChart();
                renderAdjectivePie();
            } else if (button.getAttribute('data-tab') === 'region') {
                renderRegionMap();
                renderCityChart();
            }
        });
    });
}