<template>
    <div class="admin-card">
      <h1><i class="fas fa-question-circle"></i> FAQ 管理</h1>
  
      <div class="form-container">
        <el-button type="primary" @click="openAddFAQDialog" size="large" class="submit-btn">
          <i class="fas fa-plus"></i> 添加 FAQ
        </el-button>
      </div>
  
      <el-table :data="faqList" style="width: 100%; margin-top: 20px;" class="announcement-table" stripe border>
        <el-table-column prop="question" label="问题" width="220" header-align="center" />
        <el-table-column prop="answer" label="回答" header-align="center" show-overflow-tooltip />
        <el-table-column label="操作" width="180" header-align="center" align="center">
          <template #default="scope">
            <el-button size="large" type="danger" @click="deleteFAQ(scope.row.ID)" class="delete-btn">
              <i class="fas fa-trash-alt"></i> 删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
  
      <!-- 添加 FAQ 对话框 -->
      <el-dialog v-model="addFAQDialogVisible" title="添加 FAQ">
        <el-form :model="newFAQ" label-width="80px" ref="faqFormRef">
          <el-form-item label="问题" prop="question" :rules="[{ required: true, message: '问题不能为空', trigger: 'blur' }]">
            <el-input v-model="newFAQ.question" placeholder="请输入问题" size="large" clearable />
          </el-form-item>
          <el-form-item label="回答" prop="answer" :rules="[{ required: true, message: '回答不能为空', trigger: 'blur' }]">
            <el-input v-model="newFAQ.answer" placeholder="请输入回答" size="large" clearable />
          </el-form-item>
        </el-form>
        <template #footer>
          <div class="dialog-footer">
            <el-button @click="addFAQDialogVisible = false">取 消</el-button>
            <el-button type="primary" @click="addFAQ">确 定</el-button>
          </div>
        </template>
      </el-dialog>
    </div>
  </template>
  
  <script>
  import { ref, onMounted } from 'vue';
  import { ElMessage, ElMessageBox } from 'element-plus';
  import request from '@/utils/request';
  
  export default {
    name: 'FAQManagement',
    setup() {
      const faqList = ref([]);
      const loading = ref(false);
      const addFAQDialogVisible = ref(false);
      const newFAQ = ref({ question: '', answer: '' });
      const faqFormRef = ref(null);
  
      // 获取 FAQ 列表
      const fetchFAQList = async () => {
        loading.value = true;
        try {
          const response = await request.get('/faq/list');
          faqList.value = response.data.data;
        } catch (error) {
          ElMessage.error('获取 FAQ 列表失败');
        } finally {
          loading.value = false;
        }
      };
  
      // 打开添加 FAQ 对话框
      const openAddFAQDialog = () => {
        newFAQ.value = { question: '', answer: '' }; // 重置表单
        addFAQDialogVisible.value = true;
      };
  
      // 添加 FAQ
      const addFAQ = async () => {
        if (!faqFormRef.value) return;
        await faqFormRef.value.validate(async (valid) => {
          if (!valid) return;
          try {
            await request.post('/admin/faq/add', newFAQ.value);
            ElMessage.success('FAQ 添加成功');
            addFAQDialogVisible.value = false;
            fetchFAQList();
          } catch (error) {
            ElMessage.error('添加 FAQ 失败');
          }
        });
      };
  
      // 删除 FAQ
      const deleteFAQ = (faqID) => {
        ElMessageBox.confirm('确定删除该 FAQ 吗?', '警告', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }).then(async () => {
          try {
            await request.delete(`/admin/faq/delete/${faqID}`);
            ElMessage.success('FAQ 删除成功');
            fetchFAQList();
          } catch (error) {
            ElMessage.error('删除 FAQ 失败');
          }
        });
      };
  
      onMounted(fetchFAQList);
  
      return {
        faqList,
        loading,
        addFAQDialogVisible,
        newFAQ,
        faqFormRef,
        fetchFAQList,
        openAddFAQDialog,
        addFAQ,
        deleteFAQ,
      };
    },
  };
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
  
  .create-time {
    font-family: monospace;
    color: #4a5568;
  }
  </style>
  
