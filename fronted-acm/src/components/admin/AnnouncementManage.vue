<template>
  <div class="admin-card">
    <h1><i class="fas fa-bullhorn"></i>公告管理</h1>

    <div class="form-container">
      <el-form :model="form" class="form-vertical">
        <el-form-item label="标题" class="form-item">
          <el-input v-model="form.title" placeholder="请输入公告标题" size="large" clearable />
        </el-form-item>
        <el-form-item label="内容（支持Markdown）" class="form-item">
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

export default {
  name: 'AnnouncementManage',
  data() {
    return {
      announcements: [],
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
  },
  beforeUnmount() {
    if (this.vditor) {
      this.vditor.destroy()
    }
  },
  methods: {
    initVditor() {
      this.vditor = new Vditor('vditor', {
        height: 500,
        placeholder: '请输入公告内容（支持Markdown格式）...',
        theme: 'classic',
        mode: 'wysiwyg', // 所见即所得模式
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
          '|',
          'undo',
          'redo',
          '|',
          'edit-mode',
          'preview',
          'fullscreen'
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
    }
  }
}
</script>

<style scoped>
.admin-card {
  padding: 24px;
  background: #fff;
  border-radius: 8px;
}

h1 {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 24px;
}

h1 i {
  margin-right: 8px;
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
}

:deep(.vditor-toolbar) {
  background-color: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
}

:deep(.vditor-content) {
  background-color: #fff;
}
</style>
