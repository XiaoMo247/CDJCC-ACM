<template>
  <div class="admin-card">
    <h1><i class="fas fa-bullhorn"></i>公告管理</h1>

    <div class="form-container">
      <el-form :model="form" class="form-vertical">
        <el-form-item label="标题" class="form-item">
          <el-input v-model="form.title" placeholder="请输入公告标题" size="large" clearable />
        </el-form-item>
        <el-form-item label="内容" class="form-item">
          <quill-editor
            v-model:content="form.content"
            placeholder="请输入公告内容"
            theme="snow"
            contentType="html"
            style="min-height: 300px"
            :options="editorOptions"
          />
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
      <el-table-column prop="title" label="标题" width="200" header-align="center" />
      <el-table-column prop="content" label="内容" header-align="center" show-overflow-tooltip />
      <el-table-column prop="created_at" label="创建时间" width="200" header-align="center" align="center">
        <template #default="scope">
          <span class="create-time">
            {{ formatDate(scope.row.created_at) }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" header-align="center" align="center">
        <template #default="scope">
          <el-button size="large" @click="startEdit(scope.row)" class="edit-btn">
            <i class="fas fa-edit"></i> 编辑
          </el-button>
          <el-button size="large" type="danger" @click="deleteAnnouncement(scope.row.id)" class="delete-btn">
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
import { QuillEditor } from '@vueup/vue-quill'
import '@vueup/vue-quill/dist/vue-quill.snow.css'

export default {
  name: 'AnnouncementManage',
  components: {
    QuillEditor
  },
  data() {
    return {
      announcements: [],
      form: {
        title: '',
        content: ''
      },
      editingId: null,
      editorOptions: {
        modules: {
          toolbar: [
            ['bold', 'italic', 'underline', 'strike'],
            ['blockquote', 'code-block'],
            [{ 'header': 1 }, { 'header': 2 }],
            [{ 'list': 'ordered'}, { 'list': 'bullet' }],
            [{ 'script': 'sub'}, { 'script': 'super' }],
            [{ 'indent': '-1'}, { 'indent': '+1' }],
            [{ 'direction': 'rtl' }],
            [{ 'size': ['small', false, 'large', 'huge'] }],
            [{ 'header': [1, 2, 3, 4, 5, 6, false] }],
            [{ 'color': [] }, { 'background': [] }],
            [{ 'font': [] }],
            [{ 'align': [] }],
            ['clean'],
            ['link', 'image', 'video']
          ]
        }
      }
    }
  },
  mounted() {
    this.fetchAnnouncements()
  },
  methods: {
    async fetchAnnouncements() {
      try {
        const res = await request.get('/announcement/list')
        this.announcements = res.data.announcements || []
      } catch (err) {
        ElMessage.error('获取公告列表失败')
      }
    },
    async submitAnnouncement() {
      if (!this.form.title || !this.form.content) {
        ElMessage.warning('请填写标题和内容')
        return
      }

      try {
        if (this.editingId) {
          await request.put(`/announcement/update/${this.editingId}`, this.form)
          ElMessage.success('公告更新成功')
        } else {
          await request.post('/announcement/create', this.form)
          ElMessage.success('公告创建成功')
        }
        this.resetForm()
        this.fetchAnnouncements()
      } catch (err) {
        ElMessage.error(err.response?.data?.message || '操作失败')
      }
    },
    startEdit(row) {
      this.form.title = row.title
      this.form.content = row.content
      this.editingId = row.id
    },
    cancelEdit() {
      this.resetForm()
    },
    resetForm() {
      this.form.title = ''
      this.form.content = ''
      this.editingId = null
    },
    async deleteAnnouncement(id) {
      try {
        await ElMessageBox.confirm('确定要删除该公告吗？此操作不可恢复！', '警告', {
          confirmButtonText: '确定删除',
          cancelButtonText: '取消',
          type: 'warning',
          confirmButtonClass: 'confirm-delete-btn',
          cancelButtonClass: 'cancel-delete-btn'
        })
        await request.delete(`/announcement/delete/${id}`)
        ElMessage.success('删除成功')
        this.fetchAnnouncements()
      } catch (err) {
        if (err !== 'cancel') {
          ElMessage.error(err.response?.data?.message || '删除失败')
        }
      }
    },
    formatDate(date) {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    }
  }
}
</script>

<style scoped>
.admin-card {
  background: linear-gradient(135deg, #e0f2fe, #c7d2fe);
  padding: 25px;
  border-radius: 15px;
  box-shadow: 0 6px 15px rgba(0, 0, 0, 0.08);
  color: #1e293b;
  display: flex;
  flex-direction: column;
  gap: 25px;
}

.admin-card h1 {
  font-size: 24px;
  font-weight: bold;
  margin: 0 0 15px 0;
  display: flex;
  align-items: center;
  color: #1a365d;
}

.admin-card h1 i {
  margin-right: 20px;
  width: 20px;
  height: 20px;
  color: #4a5568;
}

.form-container {
  background-color: rgba(255, 255, 255, 0.7);
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  max-width: 100%;
}

.form-vertical {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-item :deep(.el-form-item__label) {
  font-size: 16px;
  font-weight: 500;
  color: #2d3748;
  padding-bottom: 8px;
}

.form-button-group {
  display: flex;
  gap: 15px;
  margin-top: 10px;
}

.submit-btn {
  margin-top: 40px;
  font-size: 16px;
  padding: 12px 25px;
  border-radius: 8px;
  transition: all 0.3s;
}

.submit-btn i {
  margin-right: 8px;
}

.cancel-btn {
  margin-top: 40px;
  font-size: 16px;
  padding: 12px 20px;
  border-radius: 8px;
  transition: all 0.3s;
}

.cancel-btn i {
  margin-right: 8px;
}

.announcement-table {
  font-size: 16px;
  border-radius: 10px;
  overflow: hidden;
}

.announcement-table :deep(th) {
  background-color: #f8fafc;
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
}

.announcement-table :deep(td) {
  padding: 12px 0;
}

.create-time {
  font-family: monospace;
  color: #4a5568;
}

.edit-btn,
.delete-btn {
  font-size: 15px;
  padding: 10px 15px;
  border-radius: 8px;
  transition: all 0.3s;
}

.edit-btn i,
.delete-btn i {
  margin-right: 8px;
}

/* 对话框按钮样式 */
.confirm-delete-btn {
  background-color: #f56c6c;
  border-color: #f56c6c;
}

.confirm-delete-btn:hover {
  background-color: #e64c4c;
  border-color: #e64c4c;
}

.cancel-delete-btn {
  color: #909399;
}

.cancel-delete-btn:hover {
  color: #606266;
}

.form-item :deep(.ql-editor),
.form-item :deep(.ql-toolbar),
.form-item :deep(.ql-container) {
  width: 100% !important;
  box-sizing: border-box;
}

/* 修复编辑器边框和圆角 */
.form-item :deep(.ql-toolbar) {
  border-radius: 6px 6px 0 0;
}
.form-item :deep(.ql-container) {
  border-radius: 0 0 6px 6px;
  min-height: 250px; /* 你原来设置 300px 可以写这里 */
}

</style>
