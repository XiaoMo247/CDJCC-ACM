<template>
  <div class="admin-card">
    <h1><i class="fas fa-folder-open"></i>课件资源管理器</h1>

    <!-- 文件夹管理 -->
    <div class="form-container">
      <el-form :model="form" inline class="form-inline">
        <el-form-item label="文件夹名称" class="form-item">
          <el-input v-model="form.name" placeholder="请输入新文件夹名称" size="large" clearable />
        </el-form-item>
        <el-button type="primary" @click="addFolder" size="large" class="add-btn">
          <i class="fas fa-folder-plus"></i> 新建文件夹
        </el-button>
      </el-form>
    </div>

    <!-- 文件夹列表 -->
    <div class="folder-explorer">
      <el-table :data="folderList" style="width: 100%; margin-top: 20px;" class="admin-table" stripe>
        <el-table-column prop="name" label="文件夹名称" header-align="center" align="left">
          <template #default="scope">
            <div class="folder-item">
              <i class="fas fa-folder" style="color: #FFD700; margin-right: 10px;"></i>
              <span class="folder-name">{{ scope.row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" header-align="center" align="center" width="220">
          <template #default="scope">
            <span class="create-time">
              <i class="far fa-calendar-alt" style="margin-right: 5px;"></i>
              {{ formatDate(scope.row.created_at) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" header-align="center" align="center">
          <template #default="scope">
            <el-button size="large" type="primary" @click="openFileDialog(scope.row.id)" plain>
              <i class="fas fa-file-alt"></i> 查看内容
            </el-button>
            <el-button size="large" type="danger" @click="deleteFolder(scope.row.id)" plain>
              <i class="fas fa-trash-alt"></i> 删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 文件列表弹窗 -->
    <el-dialog v-model="dialogVisible" title="文件管理" width="80%" top="5vh">
      <div class="file-manager">
        <div class="file-upload-container">
          <el-upload class="upload-demo" :action="uploadUrl" :headers="uploadHeaders"
            :data="{ folder_id: currentFolderId }" :before-upload="beforeUpload" :on-success="handleUploadSuccess"
            :on-error="handleUploadError" :show-file-list="false">
            <el-button type="primary" size="large">
              <i class="fas fa-cloud-upload-alt"></i> 上传文件
            </el-button>
          </el-upload>
        </div>

        <el-table :data="fileList" style="width: 100%; margin-top: 20px;" stripe>
          <el-table-column prop="name" label="文件名" header-align="center" align="left">
            <template #default="scope">
              <div class="file-item">
                <i :class="getFileIcon(scope.row.name)" style="margin-right: 10px;"></i>
                <span class="file-name">{{ scope.row.name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="size" label="文件大小" header-align="center" align="center" width="150">
            <template #default="scope">
              <span class="file-size">
                <i class="fas fa-weight-hanging" style="margin-right: 5px;"></i>
                {{ formatFileSize(scope.row.size) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="上传时间" header-align="center" align="center" width="220">
            <template #default="scope">
              <span class="upload-time">
                <i class="far fa-clock" style="margin-right: 5px;"></i>
                {{ formatDate(scope.row.created_at) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="操作" header-align="center" align="center" width="220">
            <template #default="scope">
              <el-button size="large" type="success" @click="downloadFile(scope.row.id)" plain>
                <i class="fas fa-file-download"></i> 下载
              </el-button>
              <el-button size="large" type="danger" @click="deleteFile(scope.row.id)" plain>
                <i class="fas fa-trash"></i> 删除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import request from '@/utils/request'
import { getToken } from '@/utils/tokenManager'

export default {
  name: 'FolderManage',
  data() {
    return {
      folderList: [],
      fileList: [],
      form: {
        name: ''
      },
      dialogVisible: false,
      currentFolderId: null
    }
  },
  computed: {
    uploadHeaders() {
      const token = getToken()
      if (!token) return {}
      return { Authorization: `Bearer ${token}` }
    },
    uploadUrl() {
      return `${request.defaults.baseURL}/admin/folder/upload`
    }
  },
  mounted() {
    this.fetchFolders()
  },
  methods: {
    getFileIcon(filename) {
      const extension = filename.split('.').pop().toLowerCase()
      switch (extension) {
        case 'pdf':
          return 'fas fa-file-pdf'
        case 'doc':
        case 'docx':
          return 'fas fa-file-word'
        case 'xls':
        case 'xlsx':
          return 'fas fa-file-excel'
        case 'ppt':
        case 'pptx':
          return 'fas fa-file-powerpoint'
        case 'jpg':
        case 'jpeg':
        case 'png':
        case 'gif':
          return 'fas fa-file-image'
        case 'zip':
        case 'rar':
          return 'fas fa-file-archive'
        default:
          return 'fas fa-file'
      }
    },
    async fetchFolders() {
      try {
        const res = await request.get('/folder/list')
        this.folderList = res.data.folders || []
      } catch (error) {
        ElMessage.error('获取文件夹列表失败')
        console.error(error)
      }
    },
    async addFolder() {
      if (!this.form.name.trim()) {
        ElMessage.warning('请输入文件夹名称')
        return
      }
      try {
        await request.post('/admin/folder/add', { name: this.form.name })
        ElMessage.success('文件夹创建成功')
        this.form.name = ''
        this.fetchFolders()
      } catch (error) {
        ElMessage.error(error.response?.data?.message || '创建文件夹失败')
        console.error(error)
      }
    },
    async deleteFolder(id) {
      try {
        await ElMessageBox.confirm('确定要删除此文件夹吗？文件夹内的所有文件也将被删除', '警告', {
          confirmButtonText: '确定删除',
          cancelButtonText: '取消',
          type: 'warning'
        })
        await request.delete(`/admin/folder/delete/${id}`)
        ElMessage.success('文件夹删除成功')
        this.fetchFolders()
      } catch (error) {
        if (error !== 'cancel') {
          ElMessage.error(error.response?.data?.message || '删除文件夹失败')
          console.error(error)
        }
      }
    },
    async openFileDialog(folderId) {
      this.currentFolderId = folderId
      try {
        const res = await request.get(`/folder/files?folder_id=${folderId}`)
        this.fileList = res.data.files || []
        this.dialogVisible = true
      } catch (error) {
        ElMessage.error('获取文件列表失败')
        console.error(error)
      }
    },
    async downloadFile(fileId) {
      try {
        const res = await request.get(`/folder/download/${fileId}`)

        const { name, type, data } = res.data
        const byteCharacters = atob(data)
        const byteNumbers = new Array(byteCharacters.length).fill(0).map((_, i) => byteCharacters.charCodeAt(i))
        const byteArray = new Uint8Array(byteNumbers)
        const blob = new Blob([byteArray], { type })

        const link = document.createElement('a')
        link.href = URL.createObjectURL(blob)
        link.download = name
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)

        ElMessage.success('文件下载成功')
      } catch (error) {
        ElMessage.error('文件下载失败')
        console.error(error)
      }
    },
    async deleteFile(fileId) {
      try {
        await ElMessageBox.confirm('确定要删除此文件吗？', '警告', {
          confirmButtonText: '确定删除',
          cancelButtonText: '取消',
          type: 'warning'
        })
        await request.delete(`/admin/folder/file/${fileId}`)
        ElMessage.success('文件删除成功')
        this.openFileDialog(this.currentFolderId)
      } catch (error) {
        if (error !== 'cancel') {
          ElMessage.error(error.response?.data?.message || '文件删除失败')
          console.error(error)
        }
      }
    },
    beforeUpload(file) {
      const validTypes = [
        'application/pdf',
        'application/msword',
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
        'application/vnd.ms-powerpoint',
        'application/vnd.openxmlformats-officedocument.presentationml.presentation',
        'application/vnd.ms-excel',
        'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
        'image/jpeg',
        'image/png',
        'image/gif',
        'application/zip',
        'application/x-rar-compressed',
        'application/x-zip-compressed',
        'application/octet-stream',
        'text/markdown',
        'text/x-markdown'
      ]

      const validExtensions = ['.pdf', '.doc', '.docx', '.ppt', '.pptx', '.xls', '.xlsx', '.jpg', '.jpeg', '.png', '.gif', '.zip', '.rar', '.md']
      const fileExtension = file.name.slice(file.name.lastIndexOf('.')).toLowerCase()
      const isFileTypeValid = validTypes.includes(file.type) || validExtensions.includes(fileExtension)

      const isLt50M = file.size / 1024 / 1024 < 50
      if (!isFileTypeValid) {
        ElMessage.error('只能上传 PDF、Word、Excel、PPT、图片、压缩包或 Markdown 文件！')
      }
      if (!isLt50M) {
        ElMessage.error('文件大小不能超过 50MB！')
      }
      return isFileTypeValid && isLt50M
    }
    ,
    handleUploadSuccess(response) {
      if (response.code === 0) {
        ElMessage.success('文件上传成功')
        this.openFileDialog(this.currentFolderId)
      } else if (response.message === '上传成功') {
        ElMessage.success('文件上传成功')
        this.openFileDialog(this.currentFolderId)
      }
      else {
        ElMessage.error(response.message || '文件上传失败')
      }
    },
    handleUploadError(error) {
      ElMessage.error('文件上传失败: ' + (error.message || '未知错误'))
      console.error(error)
    },
    formatDate(date) {
      return dayjs(date).format('YYYY-MM-DD HH:mm')
    },
    formatFileSize(bytes) {
      if (bytes === 0) return '0 Bytes'
      const k = 1024
      const sizes = ['Bytes', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      const size = (bytes / Math.pow(k, i)).toFixed(2)
      return `${size} ${sizes[i]}`
    }
  }
}
</script>

<style scoped>
.form-container {
  background-color: rgba(255, 255, 255, 0.8);
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.folder-explorer {
  background-color: rgba(255, 255, 255, 0.8);
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.file-manager {
  background-color: rgba(255, 255, 255, 0.8);
  padding: 20px;
  border-radius: 10px;
}

.file-upload-container {
  margin-bottom: 20px;
  display: flex;
  justify-content: flex-end;
}

.folder-item {
  display: flex;
  align-items: center;
  font-size: 16px;
  padding: 8px 0;
}

.folder-name {
  font-weight: 500;
}

.file-item {
  display: flex;
  align-items: center;
  font-size: 16px;
  padding: 8px 0;
}

.file-name {
  font-weight: 500;
}

.create-time,
.upload-time,
.file-size {
  font-size: 15px;
  color: #555;
}

.add-btn {
  margin-left: 15px;
}

.delete-btn {
  margin-left: 10px;
}

.el-table {
  border-radius: 8px;
  overflow: hidden;
  font-size: 16px;
}

.el-table :deep(th) {
  font-size: 17px;
  font-weight: bold;
  background-color: #f8fafc !important;
}

.el-table :deep(td) {
  padding: 12px 0;
}

.el-dialog {
  border-radius: 12px;
}

.el-dialog__header {
  background-color: #f8fafc;
  border-radius: 12px 12px 0 0;
  padding: 15px 20px;
}

.el-dialog__title {
  font-size: 20px;
  font-weight: bold;
  color: #1a365d;
}

.el-button {
  font-size: 16px;
  padding: 10px 15px;
}

.el-button i {
  margin-right: 8px;
}

.el-form-item__label {
  font-size: 16px;
  font-weight: bold;
}

.el-input {
  font-size: 16px;
}

.el-message-box {
  width: 450px;
}

.el-message-box__title {
  font-size: 18px;
}
</style>
