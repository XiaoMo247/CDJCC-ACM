<template>
    <div class="settings-container">
        <!-- 用户信息展示卡片 -->
        <div class="profile-card">
            <div class="profile-header">
                <div class="avatar-container">
                    <el-avatar :size="120" :src="defaultAvatar" />
                </div>
                <div class="user-info">
                    <h2 class="username">{{ studentInfo.username }}</h2>
                    <p class="student-id">UID: {{ studentInfo.student_id }}</p>
                </div>
            </div>
        </div>

        <!-- 设置卡片 -->
        <div class="settings-card">
            <h1>
                <i class="fas fa-user-cog"></i> 账户设置
            </h1>

            <div class="form-container">
                <el-tabs v-model="activeTab" class="custom-tabs">
                    <el-tab-pane label="用户名设置" name="info">
                        <el-form :model="studentInfo" label-width="120px" class="form-inline">
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

                    <el-tab-pane label="密码设置" name="password">
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

                    <el-tab-pane label="OJ用户名设置" name="oj"> <el-form :model="ojInfo" label-width="120px" class="form-inline">
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
    </div>
</template>

<script>
import { ElMessage } from 'element-plus';
import request from '@/utils/request';
import defaultAvatar from '@/assets/user.png';

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
            defaultAvatar,
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
            }
        };
    },
    created() {
        this.fetchStudentInfo();
    },
    methods: {
        async fetchStudentInfo() {
            try {
                const response = await request.get('/student/info');
                const data = response.data?.data || {};
                this.studentInfo = {
                    student_id: data.StudentID || data.student_id || data.id || '',
                    username: data.Username || data.username || data.name || '',
                };
                this.ojInfo = {
                    cf_name: data.CfName || data.cf_name || '',
                    at_name: data.AtName || data.at_name || '',
                    nc_id: data.NcID || data.nc_id || '',
                };
            } catch (error) {
                ElMessage.error('获取用户信息失败: ' + (error.response?.data?.msg || error.message));
            }
        },

        async updateUsername() {
            try {
                await request.post('/student/update-username', { username: this.studentInfo.username });
                ElMessage.success('用户名更新成功');
            } catch (error) {
                ElMessage.error('更新用户名失败: ' + (error.response?.data?.msg || error.message));
            }
        },

        async updatePassword() {
            try {
                await this.$refs.passwordForm.validate();
                await request.post('/student/update-password', this.passwordForm);
                ElMessage.success('密码更新成功，请重新登录');
                this.$refs.passwordForm.resetFields();
            } catch (error) {
                if (error.name !== 'ValidationError') {
                    ElMessage.error('更新密码失败: ' + (error.response?.data?.msg || error.message));
                }
            }
        },

        async updateOJInfo() {
            try {
                await request.post('/student/update-info', this.ojInfo);
                ElMessage.success('OJ账号信息更新成功');
            } catch (error) {
                ElMessage.error('更新OJ账号信息失败: ' + (error.response?.data?.msg || error.message));
            }
        },
    },
};
</script>

<style scoped>
.settings-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 30px;
}

.profile-card {
    background: linear-gradient(135deg, #4f46e5, #7c3aed);
    padding: 30px;
    border-radius: 20px;
    color: white;
    box-shadow: 0 10px 25px rgba(79, 70, 229, 0.2);
}

.profile-header {
    display: flex;
    align-items: center;
}

.avatar-container {
    position: relative;
    margin-right: 30px;
}

.avatar-uploader {
    position: absolute;
    right: -10px;
    bottom: -10px;
}

.avatar-edit {
    width: 30px;
    height: 30px;
    background: #4f46e5;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.3s ease;
}

.avatar-edit:hover {
    transform: scale(1.1);
    background: #4338ca;
}

.avatar-edit i {
    color: white;
    font-size: 14px;
}

/* 隐藏上传组件的默认样式 */
:deep(.el-upload) {
    border: none;
    background: none;
}

:deep(.el-upload:hover) {
    border: none;
}

.user-avatar {
    width: 100%;
    height: 100%;
    border-radius: 50%;
    object-fit: cover;
    border: 4px solid rgba(255, 255, 255, 0.2);
}

.user-info {
    flex: 1;
}

.username {
    font-size: 2rem;
    margin: 0;
    font-weight: bold;
}

.student-id {
    font-size: 1.2rem;
    opacity: 0.8;
    margin: 5px 0 0;
}

.settings-card {
    background: white;
    padding: 30px;
    border-radius: 20px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.settings-card h1 {
    font-size: 24px;
    font-weight: bold;
    margin: 0 0 25px 0;
    display: flex;
    align-items: center;
    color: #1a365d;
}

.settings-card h1 i {
    color: #4a5568;
    margin-right: 15px;
}

.form-container {
    background-color: #f8fafc;
    padding: 25px;
    border-radius: 15px;
}

.custom-tabs {
    background: transparent;
}

.form-inline {
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.submit-btn {
    font-size: 16px;
    padding: 12px 25px;
    border-radius: 8px;
    width: fit-content;
    transition: all 0.3s ease;
}

.submit-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.submit-btn i {
    margin-right: 8px;
}

/* 移动端优化 */
@media (max-width: 768px) {
    .settings-container {
        padding: 15px;
    }

    .profile-header {
        flex-direction: column;
        text-align: center;
    }

    .avatar-container {
        margin: 0 0 20px 0;
    }

    .settings-card {
        padding: 20px;
    }

    .form-container {
        padding: 15px;
    }

    .submit-btn {
        width: 100%;
    }
}
</style>
