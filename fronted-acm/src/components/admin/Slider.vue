<template>
  <div class="admin-card">
    <h1><i class="fas fa-images"></i> 轮播图管理</h1>

    <div class="form-container">
      <el-button type="primary" @click="openAddSliderDialog" size="large" class="submit-btn">
        <i class="fas fa-plus"></i> 添加轮播图
      </el-button>
    </div>

    <el-table :data="sliderList" style="width: 100%; margin-top: 20px;" class="announcement-table" stripe border>
      <el-table-column prop="title" label="标题" width="180" header-align="center" />
      <el-table-column prop="content" label="内容" header-align="center" show-overflow-tooltip />
      <el-table-column label="图片" width="180" header-align="center">
        <template #default="{ row }">
          <el-image :src="row.image" style="width: 150px; height: 80px" fit="cover" :preview-src-list="[row.image]" />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120" header-align="center" align="center">
        <template #default="{ row }">
          <el-button size="large" type="danger" @click="deleteSlider(row.id)" class="delete-btn">
            <i class="fas fa-trash-alt"></i> 删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加轮播图对话框 -->
    <el-dialog v-model="addSliderDialogVisible" title="添加轮播图">
      <el-form :model="newSlider" label-width="80px" ref="sliderFormRef">
        <el-form-item label="标题" prop="title" :rules="[{ required: true, message: '标题不能为空', trigger: 'blur' }]">
          <el-input v-model="newSlider.title" placeholder="请输入标题" size="large" clearable />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input v-model="newSlider.content" placeholder="请输入内容" size="large" clearable />
        </el-form-item>
        <el-form-item label="图片" prop="image" :rules="[{ required: true, message: '请上传图片', trigger: 'change' }]">
          <el-upload
            class="upload-demo"
            :auto-upload="false"
            :on-change="handleFileChange"
            :show-file-list="false">
            <el-button type="primary">点击上传</el-button>
            <div v-if="newSlider.image" class="el-upload__tip">
              <el-image :src="newSlider.image" style="width: 100px; height: 60px; margin-top: 10px" />
            </div>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="addSliderDialogVisible = false">取 消</el-button>
          <el-button type="primary" @click="addSlider">确 定</el-button>
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
  name: 'SliderManagement',
  setup() {
    const sliderList = ref([]);
    const loading = ref(false);
    const addSliderDialogVisible = ref(false);
    const newSlider = ref({ title: '', content: '', image: '' });
    const sliderFormRef = ref(null);

    // 获取轮播图列表
    const fetchSliderList = async () => {
      loading.value = true;
      try {
        const response = await request.get('/slider/list');
        sliderList.value = response.data.data;
      } catch (error) {
        ElMessage.error('获取轮播图列表失败');
      } finally {
        loading.value = false;
      }
    };

    // 打开添加对话框
    const openAddSliderDialog = () => {
      newSlider.value = { title: '', content: '', image: '' };
      addSliderDialogVisible.value = true;
    };

    // 文件选择处理
    const handleFileChange = (file) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        newSlider.value.image = e.target.result;
        newSlider.value.uploadFile = file.raw; // 保存文件对象用于后续上传
      };
      reader.readAsDataURL(file.raw);
    };

    // 添加轮播图
    const addSlider = async () => {
      if (!sliderFormRef.value) return;
      await sliderFormRef.value.validate(async (valid) => {
        if (!valid) return;
        try {
          const formData = new FormData();
          formData.append('title', newSlider.value.title);
          formData.append('content', newSlider.value.content);
          if (newSlider.value.uploadFile) {
            formData.append('image', newSlider.value.uploadFile);
          }
          
          await request.post('/admin/slider/add', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
          });
          
          ElMessage.success('轮播图添加成功');
          addSliderDialogVisible.value = false;
          fetchSliderList();
        } catch (error) {
          ElMessage.error('添加轮播图失败');
        }
      });
    };

    // 删除轮播图
    const deleteSlider = (id) => {
      ElMessageBox.confirm('确定删除该轮播图吗?', '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        try {
          await request.delete(`/admin/slider/delete/${id}`);
          ElMessage.success('轮播图删除成功');
          fetchSliderList();
        } catch (error) {
          ElMessage.error('删除轮播图失败');
        }
      });
    };

    onMounted(fetchSliderList);

    return {
      sliderList,
      loading,
      addSliderDialogVisible,
      newSlider,
      sliderFormRef,
      fetchSliderList,
      openAddSliderDialog,
      handleFileChange,
      addSlider,
      deleteSlider,
    };
  },
};
</script>

<style scoped>
.admin-card {
  background: linear-gradient(135deg, #f0fdf4, #dcfce7);
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
  color: #166534;
}

.admin-card h1 i {
  margin-right: 20px;
  width: 20px;
  height: 20px;
}

.upload-demo {
  width: 100%;
}
</style>