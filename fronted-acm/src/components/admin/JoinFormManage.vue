<template>
  <div class="admin-card">
    <h1><i class="fas fa-users"></i>加入申请表管理</h1>

    <!-- 通用申请表格组件 -->
    <template v-for="(data, idx) in tableGroups" :key="idx">
      <h2>{{ data.title }}</h2>
      <el-table 
        :data="data.items" 
        style="width: 100%; margin-top: 20px;" 
        class="apply-table"
        :row-class-name="'mobile-table-row'"
      >
        <el-table-column prop="name" label="姓名" width="120" />
        <el-table-column prop="student_number" label="学号" width="150" />
        <el-table-column prop="qq_number" label="QQ号" width="150" />
        <el-table-column prop="apply" label="申请职位" width="150">
          <template #default="scope">
            {{ applyText(scope.row.apply) }}
          </template>
        </el-table-column>
        <el-table-column prop="text" label="详情信息" min-width="200" show-overflow-tooltip />
        <el-table-column prop="state" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="stateTagType(scope.row.state)" effect="light">
              {{ stateText(scope.row.state) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" :width="data.width">
          <template #default="scope">
            <div class="action-buttons">
              <button class="custom-btn detail-btn" @click="viewDetails(scope.row)">
                <i class="fas fa-eye"></i> <span class="btn-text">详情</span>
              </button>
              <button v-if="scope.row.state === 0" class="custom-btn approve-btn" @click="confirmApprove(scope.row)">
                <i class="fas fa-check"></i> <span class="btn-text">批准</span>
              </button>
              <button v-if="scope.row.state === 0" class="custom-btn reject-btn" @click="confirmReject(scope.row)">
                <i class="fas fa-times"></i> <span class="btn-text">拒绝</span>
              </button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </template>

    <!-- 申请详情弹窗 -->
    <el-dialog v-model="dialogVisible" :title="`${selectedApply.name}的申请详情`" width="90%" center>
      <template #default>
        <div class="detail-content">
          <p><strong>姓名：</strong>{{ selectedApply.name }}</p>
          <p><strong>学号：</strong>{{ selectedApply.student_number }}</p>
          <p><strong>QQ号：</strong>{{ selectedApply.qq_number }}</p>
          <p><strong>申请职位：</strong>{{ applyText(selectedApply.apply) }}</p>
          <p><strong>状态：</strong>
            <el-tag :type="stateTagType(selectedApply.state)" effect="light">
              {{ stateText(selectedApply.state) }}
            </el-tag>
          </p>
          <p><strong>自我介绍：</strong></p>
          <el-card shadow="never" class="intro-card">
            {{ selectedApply.text }}
          </el-card>
        </div>
      </template>
      <template #footer>
        <button class="custom-btn close-btn" @click="dialogVisible = false">关闭</button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import request from '@/utils/request'
import { ElMessage, ElMessageBox } from 'element-plus'

export default {
  name: 'JoinFormManage',
  data() {
    return {
      joinApplies: [],
      dialogVisible: false,
      selectedApply: {}
    }
  },
  computed: {
    pendingApplies() {
      return this.joinApplies.filter(item => item.state === 0)
    },
    approvedApplies() {
      return this.joinApplies.filter(item => item.state === 2)
    },
    rejectedApplies() {
      return this.joinApplies.filter(item => item.state === 1)
    },
    tableGroups() {
      return [
        { title: '待审核申请', items: this.pendingApplies, width: 300 },
        { title: '已批准申请', items: this.approvedApplies, width: 150 },
        { title: '已拒绝申请', items: this.rejectedApplies, width: 150 }
      ]
    }
  },
  mounted() {
    this.fetchJoinApplies()
  },
  methods: {
    applyText(code) {
      switch (code) {
        case 0: return '队员'
        case 1: return '组织委员'
        case 2: return '宣传委员'
        case 3: return '网站运维'
        default: return '未知职位'
      }
    },
    stateText(state) {
      switch (state) {
        case 0: return '待审核'
        case 1: return '已拒绝'
        case 2: return '已批准'
        default: return '未知状态'
      }
    },
    stateTagType(state) {
      switch (state) {
        case 0: return 'warning'
        case 1: return 'danger'
        case 2: return 'success'
        default: return ''
      }
    },
    async fetchJoinApplies() {
      try {
        const res = await request.get('/admin/join/get')
        this.joinApplies = res.data.applies || []
      } catch {
        ElMessage.error('获取申请数据失败')
      }
    },
    confirmApprove(row) {
      ElMessageBox.confirm(
        `确定要批准 ${row.name} (学号: ${row.student_number}) 的申请吗?`,
        '确认批准',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }
      ).then(() => {
        this.approveJoinApply(row.id)
      }).catch(() => {
        ElMessage.info('已取消操作')
      })
    },
    confirmReject(row) {
      ElMessageBox.confirm(
        `确定要拒绝 ${row.name} (学号: ${row.student_number}) 的申请吗?`,
        '确认拒绝',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }
      ).then(() => {
        this.rejectJoinApply(row.id)
      }).catch(() => {
        ElMessage.info('已取消操作')
      })
    },
    async approveJoinApply(id) {
      try {
        await request.put(`/admin/join/ac/${id}`)
        ElMessage.success('批准成功')
        this.fetchJoinApplies()
      } catch {
        ElMessage.error('批准失败')
      }
    },
    async rejectJoinApply(id) {
      try {
        await request.put(`/admin/join/wa/${id}`)
        ElMessage.success('拒绝成功')
        this.fetchJoinApplies()
      } catch {
        ElMessage.error('拒绝失败')
      }
    },
    async viewDetails(row) {
      try {
        const res = await request.get(`/admin/join/${row.id}`)
        this.selectedApply = res.data.apply
        this.dialogVisible = true
      } catch {
        ElMessage.error('获取详情失败')
      }
    }
  }
}
</script>

<style scoped>
.admin-card h2 {
  font-size: 18px;
  font-weight: bold;
  margin: 15px 0 8px 0;
  color: #2d3748;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 6px;
}

.apply-table {
  font-size: 14px;
}

.detail-content {
  font-size: 14px;
  line-height: 1.6;
}

.detail-content p {
  margin-bottom: 12px;
}

.intro-card {
  margin-top: 8px;
  padding: 12px;
  background-color: #f8fafc;
  border-radius: 8px;
  white-space: pre-line;
  line-height: 1.5;
}

.action-buttons {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.custom-btn {
  padding: 8px 10px;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.3s;
  border: none;
  display: flex;
  align-items: center;
  gap: 4px;
  min-height: 36px;
}

.custom-btn i {
  font-size: 12px;
}

.btn-text {
  display: inline-block;
}

.detail-btn {
  background-color: #f0f9ff;
  color: #1e88e5;
  border: 1px solid #bbdefb;
}

.detail-btn:hover {
  background-color: #e3f2fd;
}

.approve-btn {
  background-color: #e8f5e9;
  color: #43a047;
  border: 1px solid #c8e6c9;
}

.approve-btn:hover {
  background-color: #dcedc8;
}

.reject-btn {
  background-color: #ffebee;
  color: #e53935;
  border: 1px solid #ffcdd2;
}

.reject-btn:hover {
  background-color: #ffccbc;
}

.close-btn {
  background-color: #e3f2fd;
  color: #1976d2;
  padding: 8px 20px;
}

.close-btn:hover {
  background-color: #bbdefb;
}

.el-tag {
  font-size: 13px;
}

/* 移动端适配 */
@media screen and (max-width: 768px) {
  .admin-card h2 {
    font-size: 16px;
  }
  
  .action-buttons {
    flex-direction: column;
    gap: 4px;
  }
  
  .custom-btn {
    width: 100%;
    justify-content: center;
    padding: 8px;
  }
  
  .el-table.mobile-table-row {
    display: flex;
    flex-direction: column;
    margin-bottom: 15px;
    border: 1px solid #ebeef5;
    border-radius: 8px;
    overflow: hidden;
  }
  
  .el-table.mobile-table-row .el-table__cell {
    display: flex;
    justify-content: space-between;
    padding: 8px 12px;
    border-bottom: 1px solid #ebeef5;
  }
  
  .el-table.mobile-table-row .el-table__cell:before {
    content: attr(data-label);
    font-weight: bold;
    margin-right: 10px;
  }
  
  .el-table.mobile-table-row .el-table__cell:last-child {
    border-bottom: none;
  }
}
</style>
