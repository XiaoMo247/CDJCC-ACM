<template>
    <div class="admin-card">
        <h1><i class="fas fa-user-graduate"></i> 学生账号管理</h1>

        <div class="form-container">
            <el-form :model="form" inline class="form-inline">
                <el-form-item label="学号" class="form-item">
                    <el-input v-model="form.student_id" placeholder="请输入学生学号" size="large" clearable />
                </el-form-item>
                <el-button type="primary" @click="addStudent" size="large" class="add-btn">
                    <i class="fas fa-plus"></i> 添加学生
                </el-button>
                <el-button type="success" @click="updateAllRatings" size="large" :loading="isUpdating" class="add-btn">
                    <i class="fas fa-sync-alt"></i> 更新全部成员评分
                </el-button>
            </el-form>
        </div>

        <el-table :data="studentList" style="width: 100%; margin-top: 20px;" class="admin-table" stripe border>
            <el-table-column prop="StudentID" label="学号" header-align="center" align="center" />
            <el-table-column prop="Username" label="用户名" header-align="center" align="center" />
            <el-table-column prop="CfName" label="Codeforces 用户名" header-align="center" align="center" />
            <el-table-column prop="CfRating" label="Codeforces Rating" header-align="center" align="center" />
            <el-table-column prop="AtName" label="AtCoder 用户名" header-align="center" align="center" />
            <el-table-column prop="AtRating" label="AtCoder Rating" header-align="center" align="center" />
            <el-table-column prop="NcID" label="Nowcoder 用户 ID" header-align="center" align="center" />
            <el-table-column prop="NcRating" label="Nowcoder Rating" header-align="center" align="center" />
            <el-table-column prop="Rating" label="训练赛分数" header-align="center" align="center" />
            <el-table-column prop="Count" label="考勤次数" header-align="center" align="center" />
            <el-table-column label="创建时间" header-align="center" align="center">
                <template #default="scope">
                    <span class="create-time">{{ formatDate(scope.row.CreatedAt) }}</span>
                </template>
            </el-table-column>
            <el-table-column label="操作" width="200" header-align="center" align="center">
                <template #default="scope">
                    <el-button type="warning" size="small" @click="resetPassword(scope.row.StudentID)">重置密码</el-button>
                    <el-button type="danger" size="small" @click="deleteStudent(scope.row.StudentID)">删除</el-button>
                    <el-button type="primary" size="small" @click="updateRatings(scope.row)">更新评分</el-button>
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
    name: 'StudentManage',
    data() {
        return {
            studentList: [],
            form: {
                student_id: ''
            },
            isUpdating: false // 用于控制“更新全部评分”按钮 loading
        }
    },
    mounted() {
        this.fetchStudents()
    },
    methods: {
        async fetchStudents() {
            try {
                const res = await request.get('/admin/team-members')
                this.studentList = res.data.data || []
            } catch (err) {
                ElMessage.error('获取学生列表失败')
            }
        },
        async addStudent() {
            if (!this.form.student_id) {
                ElMessage.warning('请填写学号')
                return
            }
            try {
                await request.post('/admin/register-student', {
                    student_id: this.form.student_id
                })
                ElMessage.success('添加学生成功')
                this.form.student_id = ''
                this.fetchStudents()
            } catch (err) {
                ElMessage.error(err.response?.data?.message || '添加失败')
            }
        },
        async resetPassword(studentID) {
            try {
                const res = await request.post('/admin/reset-password', null, {
                    params: { student_id: studentID }
                })

                ElMessageBox.alert(`密码已重置为：${res.data.new_password}`, '重置成功', {
                    confirmButtonText: '确定',
                    type: 'success'
                })
            } catch (err) {
                ElMessage.error(err.response?.data?.message || '重置密码失败')
            }
        },
        async deleteStudent(studentID) {
            try {
                await ElMessageBox.confirm(
                    `确认删除学号为 ${studentID} 的学生账号？`,
                    '删除确认',
                    {
                        confirmButtonText: '确认',
                        cancelButtonText: '取消',
                        type: 'warning'
                    }
                )
                await request.delete('/admin/delete-student', {
                    params: { student_id: studentID }
                })
                ElMessage.success('删除成功')
                this.fetchStudents()
            } catch (err) {
                if (err !== 'cancel') {
                    ElMessage.error(err.response?.data?.message || '删除失败')
                }
            }
        },
        async updateRatings(student) {
            try {
                const { StudentID, CfName, AtName, NcID } = student
                await request.post('/admin/update-ratings', {
                    student_id: StudentID,
                    codeforces: CfName,
                    atcoder: AtName,
                    nowcoder: NcID
                })
                ElMessage.success('评分更新成功')
                this.fetchStudents()
            } catch (err) {
                ElMessage.error(err.response?.data?.message || '更新评分失败')
            }
        },
        async updateAllRatings() {
            try {
                this.isUpdating = true
                const res = await request.post(
                    '/admin/update-all-ratings',
                    {},
                    { timeout: 120000 } // 设置为 120000 毫秒 = 2 分钟
                )
                ElMessage.success(`更新完成：成功 ${res.data.successCount}，失败 ${res.data.failCount}`)
                this.fetchStudents()
            } catch (err) {
                if (err.code === 'ECONNABORTED') {
                    ElMessage.error('请求超时，请稍后再试')
                } else {
                    ElMessage.error(err.response?.data?.message || '更新全部评分失败')
                }
            } finally {
                this.isUpdating = false
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

.form-inline {
    display: flex;
    flex-wrap: wrap;
    gap: 15px;
    align-items: center;
}

.form-item {
    margin-bottom: 0;
}

.form-item :deep(.el-form-item__label) {
    font-size: 16px;
    font-weight: 500;
    color: #2d3748;
}

.add-btn {
    font-size: 16px;
    padding: 12px 20px;
    border-radius: 8px;
    transition: all 0.3s;
}

.add-btn i {
    margin-right: 8px;
}

.admin-table {
    font-size: 16px;
    border-radius: 10px;
    overflow: hidden;
}

.admin-table :deep(th) {
    background-color: #f8fafc;
    font-size: 16px;
    font-weight: 600;
    color: #2d3748;
}

.admin-table :deep(td) {
    padding: 12px 0;
}

.create-time {
    font-family: monospace;
    color: #4a5568;
}

@media (max-width: 768px) {
    .form-inline {
        flex-direction: column;
        align-items: stretch;
    }

    .form-item {
        width: 100%;
    }

    .add-btn {
        width: 100%;
        font-size: 15px;
    }

    .admin-table {
        font-size: 14px;
        overflow-x: auto;
    }

    .admin-table :deep(td),
    .admin-table :deep(th) {
        padding: 10px 5px;
    }
}
</style>
