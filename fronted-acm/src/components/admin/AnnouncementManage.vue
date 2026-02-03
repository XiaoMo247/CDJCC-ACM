<template>
  <div class="admin-card">
    <h1><i class="fas fa-bullhorn"></i>公告管理</h1>

    <!-- 统计图表区域 -->
    <div class="stats-section" style="margin-bottom: 24px;">
      <el-card shadow="hover">
        <template #header>
          <div style="display: flex; align-items: center;">
            <i class="fas fa-chart-line" style="margin-right: 8px;"></i>
            <span>公告观看统计 TOP10</span>
          </div>
        </template>
        <div class="charts-container" style="display: flex; gap: 24px;">
          <div id="announcementViewBarChart" style="flex: 1; height: 300px;"></div>
          <div id="announcementViewPieChart" style="flex: 1; height: 300px;"></div>
        </div>
      </el-card>
    </div>

    <div class="form-container">
      <el-form :model="form" class="form-vertical" label-position="top">
        <el-form-item label="标题" class="form-item">
          <el-input v-model="form.title" placeholder="请输入公告标题" size="large" clearable />
        </el-form-item>
        <el-form-item label="内容（支持Markdown）" class="form-item editor-item">
          <div id="vditor" style="border: 1px solid #dcdfe6; border-radius: 4px;"></div>
        </el-form-item>
        <div class="form-button-group">
          <el-button type="primary" @click="submitAnnouncement" size="large" class="submit-btn">
            <i :class="editingId ? 'fas fa-sync-alt' : 'fas fa-plus'"></i>
            {{ editingId ? '更新公告' : '新增公告' }}
          </el-button>
          <el-button v-if="editingId" @click="cancelEdit" size="large" class="cancel-btn">
            <i class="fas fa-times"></i> 取消编辑
          </el-button>
        </div>
      </el-form>
    </div>

    <el-table :data="announcements" style="width: 100%; margin-top: 20px;" class="announcement-table" stripe border>
      <el-table-column prop="title" label="标题" width="250" header-align="center" />
      <el-table-column prop="content" label="内容预览" header-align="center" show-overflow-tooltip>
        <template #default="scope">
          {{ getContentPreview(scope.row.content) }}
        </template>
      </el-table-column>
      <el-table-column prop="view_count" label="查看次数" width="120" header-align="center" align="center">
        <template #default="scope">
          <el-tag type="info">{{ scope.row.view_count || 0 }} 次</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" header-align="center" align="center">
        <template #default="scope">
          <span class="create-time">
            {{ formatDate(scope.row.created_at) }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" header-align="center" align="center">
        <template #default="scope">
          <el-button size="small" @click="startEdit(scope.row)" class="edit-btn">
            <i class="fas fa-edit"></i> 编辑
          </el-button>
          <el-button size="small" type="danger" @click="deleteAnnouncement(scope.row.id)" class="delete-btn">
            <i class="fas fa-trash-alt"></i> 删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import { getToken } from '@/utils/tokenManager'
import * as echarts from 'echarts'

export default {
  name: 'AnnouncementManage',
  data() {
    return {
      announcements: [],
      // 图表实例
      barChart: null,
      pieChart: null,
      form: {
        title: '',
        content: ''
      },
      editingId: null,
      vditor: null
    }
  },
  mounted() {
    this.fetchAnnouncements()
    this.initVditor()
    // 加载统计数据并初始化图表
    this.fetchAnnouncementViewStats()
  },
  beforeUnmount() {
    if (this.vditor) {
      this.vditor.destroy()
    }
    // 销毁图表实例
    if (this.barChart) {
      this.barChart.dispose()
    }
    if (this.pieChart) {
      this.pieChart.dispose()
    }
  },
  methods: {
    initVditor() {
      this.vditor = new Vditor('vditor', {
        height: 500,
        placeholder: '请输入公告内容（支持Markdown格式）...',
        theme: 'classic',
        mode: 'sv', // 分栏编辑模式（支持实时预览）
        preview: {
          delay: 100,
          mode: 'both', // 编辑器和预览同时显示
          hljs: {
            lineNumber: true, // 代码块显示行号
            style: 'github'
          }
        },
        toolbar: [
          'emoji',
          'headings',
          'bold',
          'italic',
          'strike',
          '|',
          'line',
          'quote',
          'list',
          'ordered-list',
          'check',
          '|',
          'code',
          'inline-code',
          'link',
          'table',
          '|',
          'upload',
          {
            hotkey: '',
            name: 'insert-image-url',
            tipPosition: 'n',
            tip: '插入网络图片',
            className: 'right',
            icon: '<svg><use xlink:href="#vditor-icon-image"></use></svg>',
            click: () => {
              const url = prompt('请输入图片URL:')
              if (url) {
                this.vditor.insertValue(`![图片](${url})`)
              }
            }
          },
          '|',
          'undo',
          'redo',
          '|',
          'edit-mode',
          {
            name: 'more',
            toolbar: [
              'fullscreen',
              'both',
              'preview'
            ]
          }
        ],
        upload: {
          url: `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'}/admin/image/upload`,
          max: 5 * 1024 * 1024, // 5MB
          fieldName: 'image',
          headers: {
            Authorization: `Bearer ${getToken()}`
          },
          format(files, responseText) {
            const response = JSON.parse(responseText)
            if (response.data && response.data.url) {
              const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
              const imageUrl = baseURL.replace('/api', '') + response.data.url
              return JSON.stringify({
                msg: '',
                code: 0,
                data: {
                  errFiles: [],
                  succMap: {
                    [files[0].name]: imageUrl
                  }
                }
              })
            }
            return JSON.stringify({
              msg: '上传失败',
              code: 1,
              data: {
                errFiles: [files[0].name],
                succMap: {}
              }
            })
          }
        },
        after: () => {
          console.log('Vditor initialized')
        }
      })
    },

    async fetchAnnouncements() {
      try {
        const res = await request.get('/announcement/list')
        this.announcements = res.data.announcements || []
      } catch (err) {
        ElMessage.error('获取公告列表失败')
      }
    },

    async submitAnnouncement() {
      if (!this.form.title) {
        ElMessage.warning('请填写标题')
        return
      }

      // 获取Vditor内容
      const content = this.vditor.getValue()
      if (!content) {
        ElMessage.warning('请填写内容')
        return
      }

      try {
        const data = {
          title: this.form.title,
          content: content
        }

        if (this.editingId) {
          await request.put(`/admin/announcement/update/${this.editingId}`, data)
          ElMessage.success('公告更新成功')
        } else {
          await request.post('/admin/announcement/create', data)
          ElMessage.success('公告创建成功')
        }
        this.resetForm()
        this.fetchAnnouncements()
      } catch (err) {
        ElMessage.error(err.response?.data?.msg || '操作失败')
      }
    },

    startEdit(row) {
      this.form.title = row.title
      this.vditor.setValue(row.content)
      this.editingId = row.id
      // 滚动到编辑器
      document.getElementById('vditor').scrollIntoView({ behavior: 'smooth', block: 'start' })
    },

    cancelEdit() {
      this.resetForm()
    },

    resetForm() {
      this.form.title = ''
      this.vditor.setValue('')
      this.editingId = null
    },

    async deleteAnnouncement(id) {
      try {
        await ElMessageBox.confirm('确定要删除该公告吗？此操作不可恢复！', '警告', {
          confirmButtonText: '确定删除',
          cancelButtonText: '取消',
          type: 'warning'
        })

        await request.delete(`/admin/announcement/delete/${id}`)
        ElMessage.success('删除成功')
        this.fetchAnnouncements()
      } catch (error) {
        if (error !== 'cancel') {
          ElMessage.error('删除失败')
        }
      }
    },

    formatDate(date) {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    },

    getContentPreview(content) {
      if (!content) return ''
      // 移除Markdown标记，获取纯文本预览
      const text = content
        .replace(/[#*`>\-\[\]()!]/g, '')
        .replace(/\n/g, ' ')
        .trim()
      return text.length > 50 ? text.substring(0, 50) + '...' : text
    },

    // 获取公告观看统计
    async fetchAnnouncementViewStats() {
      try {
        const res = await request.get('/admin/announcement/stats/views')
        const announcements = res.data.data || []

        // 等待DOM更新后初始化图表
        this.$nextTick(() => {
          this.initAnnouncementViewCharts(announcements)
        })
      } catch (err) {
        console.error('获取公告观看统计失败:', err)
      }
    },

    // 初始化公告观看统计图表
    initAnnouncementViewCharts(announcements) {
      if (!announcements || announcements.length === 0) {
        return
      }

      // 处理标题，超过15个字符截断
      const truncateTitle = (title) => {
        return title.length > 15 ? title.substring(0, 15) + '...' : title
      }

      const titles = announcements.map(a => truncateTitle(a.title))
      const viewCounts = announcements.map(a => a.view_count)

      // 初始化柱状图
      const barChartDom = document.getElementById('announcementViewBarChart')
      if (barChartDom) {
        this.barChart = echarts.init(barChartDom)
        const barOption = {
          title: {
            text: '观看量排行',
            left: 'center',
            textStyle: { fontSize: 14 }
          },
          tooltip: {
            trigger: 'axis',
            axisPointer: { type: 'shadow' },
            formatter: (params) => {
              const index = params[0].dataIndex
              return `${announcements[index].title}<br/>观看量: ${params[0].value}`
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
            data: titles,
            axisLabel: {
              interval: 0,
              rotate: 45,
              fontSize: 10
            }
          },
          yAxis: {
            type: 'value',
            name: '观看量'
          },
          series: [{
            data: viewCounts,
            type: 'bar',
            itemStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: '#f093fb' },
                { offset: 0.5, color: '#f5576c' },
                { offset: 1, color: '#f5576c' }
              ])
            },
            emphasis: {
              itemStyle: {
                color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                  { offset: 0, color: '#fa709a' },
                  { offset: 0.7, color: '#fee140' },
                  { offset: 1, color: '#fa709a' }
                ])
              }
            }
          }]
        }
        this.barChart.setOption(barOption)
      }

      // 初始化饼图
      const pieChartDom = document.getElementById('announcementViewPieChart')
      if (pieChartDom) {
        this.pieChart = echarts.init(pieChartDom)
        const pieData = announcements.map(a => ({
          value: a.view_count,
          name: truncateTitle(a.title),
          fullTitle: a.title
        }))

        const pieOption = {
          title: {
            text: '观看占比',
            left: 'center',
            textStyle: { fontSize: 14 }
          },
          tooltip: {
            trigger: 'item',
            formatter: (params) => {
              return `${params.data.fullTitle}<br/>观看量: ${params.value} (${params.percent}%)`
            }
          },
          legend: {
            orient: 'vertical',
            left: 'left',
            top: 'middle',
            textStyle: { fontSize: 10 }
          },
          series: [{
            type: 'pie',
            radius: '60%',
            data: pieData,
            emphasis: {
              itemStyle: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: 'rgba(0, 0, 0, 0.5)'
              }
            },
            label: {
              fontSize: 10
            }
          }]
        }
        this.pieChart.setOption(pieOption)
      }
    }
  }
}
</script>

<style scoped>
.admin-card {
  padding: 24px;
}

h1 {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 24px;
}

h1 i {
  margin-right: 10px;
  color: #409eff;
}

.form-container {
  margin-bottom: 30px;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;
}

.form-vertical {
  max-width: 100%;
}

.form-item {
  margin-bottom: 20px;
}

/* 编辑器表单项特殊样式 */
.editor-item {
  width: 100%;
}

.editor-item :deep(.el-form-item__label) {
  font-weight: 600;
  font-size: 15px;
  margin-bottom: 10px;
}

.editor-item :deep(.el-form-item__content) {
  width: 100%;
  max-width: 100%;
}

.form-button-group {
  display: flex;
  gap: 12px;
  margin-top: 20px;
}

.submit-btn, .cancel-btn {
  min-width: 120px;
}

.announcement-table {
  border-radius: 8px;
  overflow: hidden;
}

.announcement-table :deep(.el-table__header) {
  background-color: #f5f7fa;
}

.create-time {
  font-size: 13px;
  color: #909399;
}

.edit-btn, .delete-btn {
  margin: 0 4px;
}

/* Vditor样式覆盖 */
:deep(.vditor) {
  border-radius: 4px;
  width: 100%;
  max-width: 100%;
}

:deep(.vditor-toolbar) {
  background-color: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
}

:deep(.vditor-content) {
  background-color: #fff;
}

/* 确保编辑器铺满表单容器 */
#vditor {
  width: 100% !important;
  max-width: 100% !important;
  box-sizing: border-box;
}

.form-item :deep(.el-form-item__content) {
  width: 100%;
  max-width: 100%;
}
</style>
