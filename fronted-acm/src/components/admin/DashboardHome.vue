<template>
  <div class="dashboard-home">
    <h1><i class="fas fa-tachometer-alt"></i> 管理主页</h1>

    <!-- 统计卡片区域 -->
    <div class="stats-cards">
      <el-card class="stat-card card-blue" shadow="hover" @click="navigateTo('announcement')">
        <div class="card-content">
          <div class="card-icon">
            <i class="fas fa-bullhorn"></i>
          </div>
          <div class="card-info">
            <div class="card-value">{{ stats.announcementCount }}</div>
            <div class="card-label">公告总数</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card card-orange" shadow="hover" @click="navigateTo('join')">
        <div class="card-content">
          <div class="card-icon">
            <i class="fas fa-clipboard-list"></i>
          </div>
          <div class="card-info">
            <div class="card-value">{{ stats.pendingApplicationCount }}</div>
            <div class="card-label">待处理申请表</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card card-green" shadow="hover" @click="navigateTo('user')">
        <div class="card-content">
          <div class="card-icon">
            <i class="fas fa-users"></i>
          </div>
          <div class="card-info">
            <div class="card-value">{{ stats.userCount }}</div>
            <div class="card-label">用户总数</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card card-purple" shadow="hover" @click="navigateTo('student')">
        <div class="card-content">
          <div class="card-icon">
            <i class="fas fa-user-graduate"></i>
          </div>
          <div class="card-info">
            <div class="card-value">{{ stats.studentCount }}</div>
            <div class="card-label">学生总数</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card card-cyan" shadow="hover" @click="navigateTo('contest')">
        <div class="card-content">
          <div class="card-icon">
            <i class="fas fa-trophy"></i>
          </div>
          <div class="card-info">
            <div class="card-value">{{ stats.contestCount }}</div>
            <div class="card-label">比赛总数</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 图表区域 -->
    <div class="charts-container">
      <el-card class="chart-card" shadow="hover">
        <template #header>
          <div class="chart-header">
            <i class="fas fa-chart-bar"></i> 各项数据统计
          </div>
        </template>
        <div id="barChart" style="width: 100%; height: 400px;"></div>
      </el-card>

      <el-card class="chart-card" shadow="hover">
        <template #header>
          <div class="chart-header">
            <i class="fas fa-chart-pie"></i> 数据占比分析
          </div>
        </template>
        <div id="pieChart" style="width: 100%; height: 400px;"></div>
      </el-card>
    </div>
  </div>
</template>

<script>
import request from '@/utils/request'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'

export default {
  name: 'DashboardHome',
  data() {
    return {
      stats: {
        userCount: 0,
        studentCount: 0,
        contestCount: 0,
        announcementCount: 0,
        pendingApplicationCount: 0
      },
      barChart: null,
      pieChart: null
    }
  },
  mounted() {
    this.fetchStats()
  },
  beforeUnmount() {
    // 销毁图表实例，避免内存泄漏
    if (this.barChart) {
      this.barChart.dispose()
    }
    if (this.pieChart) {
      this.pieChart.dispose()
    }
  },
  methods: {
    async fetchStats() {
      try {
        const res = await request.get('/admin/dashboard/stats')
        this.stats = res.data

        // 数据获取成功后渲染图表
        this.$nextTick(() => {
          this.initBarChart()
          this.initPieChart()
        })
      } catch (error) {
        ElMessage.error('获取统计数据失败')
        console.error(error)
      }
    },

    initBarChart() {
      const chartDom = document.getElementById('barChart')
      if (!chartDom) return

      this.barChart = echarts.init(chartDom)

      const option = {
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
          data: ['公告', '待处理申请', '用户', '学生', '比赛'],
          axisLabel: {
            fontSize: 14
          }
        },
        yAxis: {
          type: 'value'
        },
        series: [
          {
            name: '数量',
            type: 'bar',
            data: [
              this.stats.announcementCount,
              this.stats.pendingApplicationCount,
              this.stats.userCount,
              this.stats.studentCount,
              this.stats.contestCount
            ],
            itemStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: '#83bff6' },
                { offset: 0.5, color: '#188df0' },
                { offset: 1, color: '#188df0' }
              ])
            },
            emphasis: {
              itemStyle: {
                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                  { offset: 0, color: '#2378f7' },
                  { offset: 0.7, color: '#2378f7' },
                  { offset: 1, color: '#83bff6' }
                ])
              }
            },
            label: {
              show: true,
              position: 'top',
              fontSize: 14,
              fontWeight: 'bold'
            }
          }
        ]
      }

      this.barChart.setOption(option)

      // 响应式调整图表大小
      window.addEventListener('resize', () => {
        this.barChart?.resize()
      })
    },

    initPieChart() {
      const chartDom = document.getElementById('pieChart')
      if (!chartDom) return

      this.pieChart = echarts.init(chartDom)

      const option = {
        tooltip: {
          trigger: 'item',
          formatter: '{a} <br/>{b}: {c} ({d}%)'
        },
        legend: {
          orient: 'vertical',
          left: 'left',
          top: 'center'
        },
        series: [
          {
            name: '数据统计',
            type: 'pie',
            radius: ['40%', '70%'],
            avoidLabelOverlap: false,
            itemStyle: {
              borderRadius: 10,
              borderColor: '#fff',
              borderWidth: 2
            },
            label: {
              show: true,
              formatter: '{b}: {c}'
            },
            emphasis: {
              label: {
                show: true,
                fontSize: 16,
                fontWeight: 'bold'
              }
            },
            labelLine: {
              show: true
            },
            data: [
              {
                value: this.stats.announcementCount,
                name: '公告',
                itemStyle: { color: '#5470c6' }
              },
              {
                value: this.stats.pendingApplicationCount,
                name: '待处理申请',
                itemStyle: { color: '#91cc75' }
              },
              {
                value: this.stats.userCount,
                name: '用户',
                itemStyle: { color: '#fac858' }
              },
              {
                value: this.stats.studentCount,
                name: '学生',
                itemStyle: { color: '#ee6666' }
              },
              {
                value: this.stats.contestCount,
                name: '比赛',
                itemStyle: { color: '#73c0de' }
              }
            ]
          }
        ]
      }

      this.pieChart.setOption(option)

      // 响应式调整图表大小
      window.addEventListener('resize', () => {
        this.pieChart?.resize()
      })
    },

    navigateTo(tab) {
      // 通过emit通知父组件切换标签页
      this.$emit('change-tab', tab)
    }
  }
}
</script>

<style scoped>
.dashboard-home {
  padding: 20px;
}

h1 {
  font-size: 28px;
  color: #333;
  margin-bottom: 30px;
  font-weight: 600;
}

h1 i {
  margin-right: 12px;
  color: #409eff;
}

/* 统计卡片区域 */
.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 12px;
  overflow: hidden;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15) !important;
}

.card-content {
  display: flex;
  align-items: center;
  padding: 10px;
}

.card-icon {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  font-size: 36px;
  margin-right: 20px;
}

.card-info {
  flex: 1;
}

.card-value {
  font-size: 36px;
  font-weight: bold;
  margin-bottom: 5px;
}

.card-label {
  font-size: 16px;
  color: #666;
  font-weight: 500;
}

/* 不同颜色的卡片 */
.card-blue .card-icon {
  background-color: rgba(64, 158, 255, 0.1);
  color: #409eff;
}

.card-blue .card-value {
  color: #409eff;
}

.card-orange .card-icon {
  background-color: rgba(255, 152, 0, 0.1);
  color: #ff9800;
}

.card-orange .card-value {
  color: #ff9800;
}

.card-green .card-icon {
  background-color: rgba(76, 175, 80, 0.1);
  color: #4caf50;
}

.card-green .card-value {
  color: #4caf50;
}

.card-purple .card-icon {
  background-color: rgba(156, 39, 176, 0.1);
  color: #9c27b0;
}

.card-purple .card-value {
  color: #9c27b0;
}

.card-cyan .card-icon {
  background-color: rgba(0, 188, 212, 0.1);
  color: #00bcd4;
}

.card-cyan .card-value {
  color: #00bcd4;
}

/* 图表区域 */
.charts-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(500px, 1fr));
  gap: 20px;
}

.chart-card {
  border-radius: 12px;
}

.chart-header {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.chart-header i {
  margin-right: 8px;
  color: #409eff;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .stats-cards {
    grid-template-columns: 1fr;
  }

  .charts-container {
    grid-template-columns: 1fr;
  }

  .card-icon {
    width: 60px;
    height: 60px;
    font-size: 28px;
  }

  .card-value {
    font-size: 28px;
  }

  .card-label {
    font-size: 14px;
  }
}
</style>
