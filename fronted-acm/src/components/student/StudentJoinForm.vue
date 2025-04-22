<template>
    <div class="join-card">
      <h1><i class="fas fa-edit"></i>加入申请表</h1>
  
      <!-- 提交申请表部分 -->
      <div class="form-container">
        <el-form :model="form" label-width="100px" class="form-inline">
          <el-form-item label="学号">
            <el-input v-model="form.student_number" placeholder="请输入学号" size="large" clearable />
          </el-form-item>
          <el-form-item label="QQ号">
            <el-input v-model="form.qq_number" placeholder="请输入QQ号" size="large" clearable />
          </el-form-item>
          <el-form-item label="申请职位">
            <el-select v-model="form.apply" placeholder="请选择职位" size="large">
              <el-option label="队员" :value="0" />
              <el-option label="组织委员" :value="1" />
              <el-option label="宣传委员" :value="2" />
              <el-option label="网站运维" :value="3" />
            </el-select>
          </el-form-item>
          <el-form-item label="真实姓名">
            <el-input v-model="form.name" placeholder="请输入真实姓名" size="large" clearable />
          </el-form-item>
          <el-form-item label="自我介绍" class="form-item">
            <el-input
              v-model="form.text"
              type="textarea"
              placeholder="介绍一下自己~"
              :rows="3"
              resize="none"
              size="large"
              class="content-textarea"
            />
          </el-form-item>
  
          <el-button type="primary" @click="submitApply" size="large" class="submit-btn">
            <i class="fas fa-paper-plane"></i> 提交申请
          </el-button>
        </el-form>
      </div>
  
      <!-- 申请记录部分 -->
      <h2><i class="fas fa-list"></i>我的申请记录</h2>
      <el-table :data="myApplies" style="width: 100%" class="apply-table" stripe border>
        <el-table-column prop="name" label="姓名" align="center" />
        <el-table-column prop="student_number" label="学号" align="center" />
        <el-table-column prop="qq_number" label="QQ号" align="center" />
        <el-table-column prop="apply" label="申请职位" align="center">
          <template #default="scope">
            <span>{{ getApplyPositionText(scope.row.apply) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="text" label="自我介绍" align="center" />
        <el-table-column label="状态" align="center">
          <template #default="scope">
            <el-tag :type="statusTagType(scope.row.state)">
              {{ statusText(scope.row.state) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="提交时间" align="center">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
    </div>
  </template>
  
  <script>
  import request from '@/utils/request'
  import { ElMessage } from 'element-plus'
  import dayjs from 'dayjs'
  
  export default {
    name: 'StudentJoinForm',
    data() {
      return {
        form: {
          student_number: '',
          qq_number: '',
          apply: '',
          name: '',
          text: ''
        },
        myApplies: []
      }
    },
    mounted() {
      this.getMyApplies()
    },
    methods: {
      async submitApply() {
        const { student_number, qq_number, apply, name, text } = this.form
        if (!student_number || !qq_number || apply === '' || !name || !text) {
          ElMessage.warning('请填写所有信息')
          return
        }
        try {
          await request.post('/join/send', this.form)
          ElMessage.success('申请提交成功')
          this.form = {
            student_number: '',
            qq_number: '',
            apply: '',
            name: '',
            text: ''
          }
          this.getMyApplies()
        } catch (err) {
          ElMessage.error(err.response?.data?.message || '提交失败')
        }
      },
      async getMyApplies() {
        try {
          const res = await request.get('/join/my')
          this.myApplies = res.data.data || []
        } catch (err) {
          ElMessage.error('获取申请记录失败')
        }
      },
      formatDate(date) {
        return dayjs(date).format('YYYY-MM-DD HH:mm')
      },
      statusText(state) {
        switch (state) {
          case 0:
            return '待审核'
          case 2:
            return '已通过'
          case 1:
            return '已拒绝'
          default:
            return '未知'
        }
      },
      statusTagType(state) {
        return state === 2 ? 'success' : state === 1 ? 'danger' : 'info'
      },
      getApplyPositionText(value) {
        const positions = {
          0: '队员',
          1: '组织委员',
          2: '宣传委员',
          3: '网站运维'
        }
        return positions[value] || '未知'
      }
    }
  }
  </script>
  
  <style scoped>
  .join-card {
    background: linear-gradient(135deg, #e0f2fe, #c7d2fe);
    padding: 25px;
    border-radius: 15px;
    box-shadow: 0 6px 15px rgba(0, 0, 0, 0.08);
    color: #1e293b;
    display: flex;
    flex-direction: column;
    gap: 25px;
  }
  
  .join-card h1,
  .join-card h2 {
    font-size: 24px;
    font-weight: bold;
    margin: 0 0 15px 0;
    display: flex;
    align-items: center;
    color: #1a365d;
  }
  
  .join-card h1 i,
  .join-card h2 i {
    color: #4a5568;
    margin-right: 20px;
    width: 20px;
    height: 20px;
  }
  
  .form-container {
    background-color: rgba(255, 255, 255, 0.8);
    padding: 20px;
    border-radius: 10px;
  }
  
  .form-inline {
    display: flex;
    flex-direction: column;
    gap: 15px;
  }
  
  .submit-btn {
    font-size: 16px;
    padding: 12px 20px;
    border-radius: 8px;
    width: fit-content;
  }
  
  .submit-btn i {
    margin-right: 8px;
  }
  
  .apply-table {
    font-size: 15px;
    border-radius: 10px;
    overflow: hidden;
  }
  </style>
  