<template>
    <div class="admin-card">
      <h1><i class="fas fa-trophy"></i>比赛管理</h1>
  
      <div class="form-container">
        <el-form :model="form" class="form-vertical">
          <el-form-item label="比赛标题" class="form-item">
            <el-input v-model="form.title" placeholder="请输入比赛标题" size="large" clearable />
          </el-form-item>
          <el-form-item label="比赛简介" class="form-item">
            <el-input v-model="form.text" placeholder="请输入比赛简介" size="large" clearable />
          </el-form-item>
          <el-form-item label="比赛链接" class="form-item">
            <el-input v-model="form.link" placeholder="请输入比赛链接" size="large" clearable />
          </el-form-item>
          <el-form-item label="比赛时间" class="form-item">
            <el-date-picker
              v-model="form.time"
              type="datetime"
              placeholder="选择比赛时间"
              size="large"
              style="width: 100%;"
            />
          </el-form-item>
  
          <div class="form-button-group">
            <el-button type="primary" @click="submitCompetition" size="large" class="submit-btn">
              <i class="fas fa-plus"></i> 新增比赛
            </el-button>
          </div>
        </el-form>
      </div>
  
      <el-table :data="competitions" style="width: 100%; margin-top: 20px;" class="announcement-table" stripe border>
        <el-table-column prop="title" label="比赛标题" width="220" header-align="center" />
        <el-table-column prop="text" label="简介" header-align="center" show-overflow-tooltip />
        <el-table-column prop="link" label="链接" header-align="center" show-overflow-tooltip>
          <template #default="scope">
            <a v-if="scope.row.link" :href="scope.row.link" target="_blank" style="color: #3b82f6;">
              {{ scope.row.link }}
            </a>
          </template>
        </el-table-column>
        <el-table-column prop="time" label="时间" width="180" header-align="center" align="center">
          <template #default="scope">
            <span class="create-time">
              {{ formatDate(scope.row.time) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" header-align="center" align="center">
          <template #default="scope">
            <el-button size="large" type="danger" @click="deleteCompetition(scope.row.id)" class="delete-btn">
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
  
  export default {
    name: 'CompetitionManage',
    data() {
      return {
        competitions: [],
        form: {
          title: '',
          text: '',
          link: '',
          time: ''
        }
      }
    },
    mounted() {
      this.fetchCompetitions()
    },
    methods: {
      async fetchCompetitions() {
        try {
          const res = await request.get('/contest/list')
          this.competitions = res.data.data || []
        } catch (err) {
          ElMessage.error('获取比赛列表失败')
        }
      },
      async submitCompetition() {
        if (!this.form.title || !this.form.time) {
          ElMessage.warning('请填写比赛标题和时间')
          return
        }
        try {
          await request.post('/admin/contest/upload', this.form)
          ElMessage.success('比赛创建成功')
          this.resetForm()
          this.fetchCompetitions()
        } catch (err) {
          ElMessage.error(err.response?.data?.message || '创建失败')
        }
      },
      resetForm() {
        this.form.title = ''
        this.form.text = ''
        this.form.link = ''
        this.form.time = ''
      },
      async deleteCompetition(id) {
        try {
          await ElMessageBox.confirm('确定要删除该比赛吗？', '警告', {
            confirmButtonText: '确定删除',
            cancelButtonText: '取消',
            type: 'warning',
            confirmButtonClass: 'confirm-delete-btn',
            cancelButtonClass: 'cancel-delete-btn'
          })
          await request.delete(`/admin/contest/${id}`)
          ElMessage.success('删除成功')
          this.fetchCompetitions()
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
  .form-container {
    background-color: rgba(255, 255, 255, 0.7);
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
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
  
  .submit-btn, .delete-btn {
    font-size: 16px;
    padding: 12px 20px;
    border-radius: 8px;
    transition: all 0.3s;
  }
  
  .submit-btn i, .delete-btn i {
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
  
  .create-time {
    font-family: monospace;
    color: #4a5568;
  }
  </style>
  
