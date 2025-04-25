<template>
    <div class="settings-card">
        <h1>
            <i class="fas fa-user-cog"></i> 账户设置
        </h1>

        <div class="form-container">
            <el-tabs v-model="activeTab">
                <el-tab-pane label="基本信息" name="info">
                    <el-form :model="studentInfo" label-width="120px" class="form-inline">
                        <el-form-item label="学号">
                            <el-input v-model="studentInfo.student_id" disabled></el-input>
                        </el-form-item>
                        <el-form-item label="用户名">
                            <el-input v-model="studentInfo.username"></el-input>
                        </el-form-item>
                        <el-form-item>
                            <el-button type="primary" class="submit-btn" @click="updateUsername">
                                <i class="fas fa-save"></i> 更新用户名
                            </el-button>
                        </el-form-item>
                    </el-form>
                </el-tab-pane>

                <el-tab-pane label="修改密码" name="password">
                    <el-form :model="passwordForm" :rules="passwordRules" ref="passwordForm" label-width="120px" class="form-inline">
                        <el-form-item label="原密码" prop="old_password">
                            <el-input v-model="passwordForm.old_password" type="password" show-password></el-input>
                        </el-form-item>
                        <el-form-item label="新密码" prop="new_password">
                            <el-input v-model="passwordForm.new_password" type="password" show-password></el-input>
                        </el-form-item>
                        <el-form-item label="确认密码" prop="confirm">
                            <el-input v-model="passwordForm.confirm" type="password" show-password></el-input>
                        </el-form-item>
                        <el-form-item>
                            <el-button type="primary" class="submit-btn" @click="updatePassword">
                                <i class="fas fa-key"></i> 修改密码
                            </el-button>
                        </el-form-item>
                    </el-form>
                </el-tab-pane>

                <el-tab-pane label="OJ账号" name="oj">
                    <el-form :model="ojInfo" label-width="120px" class="form-inline">
                        <el-form-item label="Codeforces">
                            <el-input v-model="ojInfo.cf_name" placeholder="输入Codeforces用户名"></el-input>
                        </el-form-item>
                        <el-form-item label="AtCoder">
                            <el-input v-model="ojInfo.at_name" placeholder="输入AtCoder用户名"></el-input>
                        </el-form-item>
                        <el-form-item label="Nowcoder">
                            <el-input v-model="ojInfo.nc_id" placeholder="输入Nowcoder用户ID"></el-input>
                        </el-form-item>
                        <el-form-item>
                            <el-button type="primary" class="submit-btn" @click="updateOJInfo">
                                <i class="fas fa-code"></i> 更新OJ信息
                            </el-button>
                        </el-form-item>
                    </el-form>
                </el-tab-pane>
            </el-tabs>
        </div>
    </div>
</template>
  <script>
  import { ElMessage } from 'element-plus';
  import request from '@/utils/request'; // 确保你有这个封装的request工具
  
  export default {
    data() {
      const validatePassword = (rule, value, callback) => {
        if (value !== this.passwordForm.new_password) {
          callback(new Error('两次输入密码不一致!'));
        } else {
          callback();
        }
      };
  
      return {
        activeTab: 'info',
        studentInfo: {
          student_id: '',
          username: '',
        },
        passwordForm: {
          old_password: '',
          new_password: '',
          confirm: '',
        },
        passwordRules: {
          old_password: [
            { required: true, message: '请输入原密码', trigger: 'blur' },
            { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
          ],
          new_password: [
            { required: true, message: '请输入新密码', trigger: 'blur' },
            { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
          ],
          confirm: [
            { required: true, message: '请再次输入新密码', trigger: 'blur' },
            { validator: validatePassword, trigger: 'blur' },
          ],
        },
        ojInfo: {
          cf_name: '',
          at_name: '',
          nc_id: '',
        },
      };
    },
    created() {
      this.fetchStudentInfo();
    },
    methods: {
      // 获取学生信息
      async fetchStudentInfo() {
        try {
          const response = await request.get('/student/info', {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('member_token')}`, // 假设你保存在localStorage中
            },
          });
          console.log(response.data.data)
          this.studentInfo = {
            student_id: response.data.data.StudentID,
            username: response.data.data.Username,
          };
          this.ojInfo = {
            cf_name: response.data.data.CfName || '',
            at_name: response.data.data.AtName || '',
            nc_id: response.data.data.NcID || '',
          };
        } catch (error) {
          ElMessage.error('获取用户信息失败: ' + (error.response?.data?.msg || error.message));
        }
      },
  
      // 更新用户名
      async updateUsername() {
        try {
          await request.post(
            '/student/update-username',
            { username: this.studentInfo.username },
            {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('member_token')}`, // 使用请求头传递token
              },
            }
          );
          ElMessage.success('用户名更新成功');
        } catch (error) {
          ElMessage.error('更新用户名失败: ' + (error.response?.data?.msg || error.message));
        }
      },
  
      // 更新密码
      async updatePassword() {
        try {
          await this.$refs.passwordForm.validate();
          await request.post(
            '/student/update-password',
            this.passwordForm,
            {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('member_token')}`, // 使用请求头传递token
              },
            }
          );
          ElMessage.success('密码更新成功，请重新登录');
          this.$refs.passwordForm.resetFields();
          // 这里可以添加重新登录的逻辑
        } catch (error) {
          if (error.name !== 'ValidationError') {
            ElMessage.error('更新密码失败: ' + (error.response?.data?.msg || error.message));
          }
        }
      },
  
      // 更新OJ账号信息
      async updateOJInfo() {
        try {
          await request.post(
            '/student/update-info',
            this.ojInfo,
            {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('member_token')}`, // 使用请求头传递token
              },
            }
          );
          ElMessage.success('OJ账号信息更新成功');
        } catch (error) {
          ElMessage.error('更新OJ账号信息失败: ' + (error.response?.data?.msg || error.message));
        }
      },
    },
  };
  </script>

<style scoped>
.settings-card {
    background: linear-gradient(135deg, #e0f2fe, #c7d2fe);
    padding: 25px;
    border-radius: 15px;
    box-shadow: 0 6px 15px rgba(0, 0, 0, 0.08);
    color: #1e293b;
    display: flex;
    flex-direction: column;
    gap: 25px;
}

.settings-card h1 {
    font-size: 24px;
    font-weight: bold;
    margin: 0 0 15px 0;
    display: flex;
    align-items: center;
    color: #1a365d;
}

.settings-card h1 i {
    color: #4a5568;
    margin-right: 20px;
    width: 20px;
    height: 20px;
}

.form-container {
    background-color: rgba(255, 255, 255, 0.85);
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
</style>