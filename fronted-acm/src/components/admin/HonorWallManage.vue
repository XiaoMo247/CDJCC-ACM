<template>
    <div class="admin-card">
      <h1><i class="fas fa-trophy"></i>荣誉墙管理</h1>
  
      <div class="form-container">
        <el-form :model="form" class="form-vertical">
          <el-form-item label="姓名" class="form-item">
            <el-input v-model="form.name" placeholder="请输入姓名" size="large" clearable />
          </el-form-item>
          <el-form-item label="获奖年份" class="form-item">
            <el-input v-model.number="form.year" placeholder="请输入获奖年份" size="large" clearable />
          </el-form-item>
          <el-form-item label="赛事类型" class="form-item">
            <el-select v-model="form.contest" placeholder="请选择赛事类型" size="large" style="width: 100%;">
              <el-option label="ICPC" value="ICPC" />
              <el-option label="CCPC" value="CCPC" />
              <el-option label="蓝桥杯" value="蓝桥杯" />
              <el-option label="百度之星" value="百度之星" />
            </el-select>
          </el-form-item>
          <el-form-item label="获奖级别" class="form-item">
            <el-select v-model="form.level" placeholder="请选择获奖级别" size="large" style="width: 100%;">
              <el-option label="金奖" value="金奖" />
              <el-option label="银奖" value="银奖" />
              <el-option label="铜奖" value="铜奖" />
              <el-option label="国一" value="国一" />
              <el-option label="国二" value="国二" />
              <el-option label="国三" value="国三" />
              <el-option label="省一" value="省一" />
              <el-option label="省二" value="省二" />
              <el-option label="省三" value="省三" />
            </el-select>
          </el-form-item>
          <el-form-item label="获奖时年级" class="form-item">
            <el-select v-model="form.grade" placeholder="请选择年级" size="large" style="width: 100%;">
              <el-option label="大一" :value="1" />
              <el-option label="大二" :value="2" />
              <el-option label="大三" :value="3" />
              <el-option label="大四" :value="4" />
            </el-select>
          </el-form-item>
          <el-form-item label="专业" class="form-item">
            <el-select v-model="form.major" placeholder="请选择专业" size="large" style="width: 100%;">
              <el-option label="软件工程" value="软件工程" />
              <el-option label="计算机科学与技术" value="计算机科学与技术" />
              <el-option label="人工智能" value="人工智能" />
              <el-option label="数字媒体技术" value="数字媒体技术" />
            </el-select>
          </el-form-item>
  
          <div class="form-button-group">
            <el-button type="primary" @click="submitHonor" size="large" class="submit-btn">
              <i class="fas fa-plus"></i> 新增荣誉
            </el-button>
          </div>
        </el-form>
      </div>
  
      <el-table :data="honors" style="width: 100%; margin-top: 20px;" class="announcement-table" stripe border>
        <el-table-column prop="name" label="姓名" width="120" header-align="center" align="center" />
        <el-table-column prop="year" label="年份" width="100" header-align="center" align="center" />
        <el-table-column prop="contest" label="赛事类型" width="120" header-align="center" align="center" />
        <el-table-column prop="level" label="获奖级别" width="120" header-align="center" align="center" />
        <el-table-column prop="grade" label="年级" width="100" header-align="center" align="center">
          <template #default="scope">
            <span>{{ getGradeText(scope.row.grade) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="major" label="专业" header-align="center" show-overflow-tooltip />
        <el-table-column label="操作" width="180" header-align="center" align="center">
          <template #default="scope">
            <el-button size="large" type="danger" @click="deleteHonor(scope.row.id)" class="delete-btn">
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
  
  export default {
    name: 'HonorWallManage',
    data() {
      return {
        honors: [],
        form: {
          name: '',
          year: null,
          contest: '',
          level: '',
          grade: null,
          major: ''
        }
      }
    },
    mounted() {
      this.fetchHonors()
    },
    methods: {
      async fetchHonors() {
        try {
          const res = await request.get('/honor')
          this.honors = res.data.data || []
        } catch (err) {
          ElMessage.error('获取荣誉列表失败')
        }
      },
      async submitHonor() {
        // 验证年份是否为有效数字
        if (typeof this.form.year !== 'number' || isNaN(this.form.year)) {
          ElMessage.warning('年份必须是有效的数字')
          return
        }
        
        if (!this.form.name || !this.form.contest || !this.form.level || !this.form.grade) {
          ElMessage.warning('请填写必填项')
          return
        }
        
        try {
          await request.post('/admin/honor', this.form)
          ElMessage.success('荣誉记录添加成功')
          this.resetForm()
          this.fetchHonors()
        } catch (err) {
          ElMessage.error(err.response?.data?.message || '添加失败')
        }
      },
      resetForm() {
        this.form = {
          name: '',
          year: null,
          contest: '',
          level: '',
          grade: null,
          major: ''
        }
      },
      async deleteHonor(id) {
        try {
          await ElMessageBox.confirm('确定要删除该荣誉记录吗？', '警告', {
            confirmButtonText: '确定删除',
            cancelButtonText: '取消',
            type: 'warning',
            confirmButtonClass: 'confirm-delete-btn',
            cancelButtonClass: 'cancel-delete-btn'
          })
          await request.delete(`/admin/honor/${id}`)
          ElMessage.success('删除成功')
          this.fetchHonors()
        } catch (err) {
          if (err !== 'cancel') {
            ElMessage.error(err.response?.data?.message || '删除失败')
          }
        }
      },
      getGradeText(grade) {
        const grades = ['', '大一', '大二', '大三', '大四']
        return grades[grade] || '未知'
      }
    }
  }
  </script>
  
  <style scoped>
  h1 {
    font-size: 24px;
    color: #303133;
    margin-bottom: 24px;
    font-weight: 600;
  }

  h1 i {
    margin-right: 10px;
    color: #409eff;
  }

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
  </style>
  
