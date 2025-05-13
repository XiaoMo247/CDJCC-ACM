<template>
    <div class="competition-page">
        <!-- 页面标题 -->
        <div class="page-header">
            <h1> <span style="font-size: 3rem;">🏆</span><span class="page-title">比赛信息中心</span></h1>
            <p class="page-subtitle">了解近期重要比赛，助力你的成长之路！</p>
        </div>

        <!-- 比赛列表容器 -->
        <div class="competition-container">
            <!-- 即将开始的比赛 -->
            <div class="section">
                <h2 class="section-title">
                    <el-icon>
                        <AlarmClock />
                    </el-icon> 即将开始
                    <span class="badge">{{ upcomingContests.length }}</span>
                </h2>

                <div class="competition-list">
                    <el-skeleton :loading="loading" animated :count="3">
                        <template #template>
                            <div class="competition-item skeleton">
                                <el-skeleton-item variant="image"
                                    style="width: 100%; height: 120px; border-radius: 12px 12px 0 0;" />
                                <div style="padding: 16px;">
                                    <el-skeleton-item variant="h3" style="width: 60%; margin-bottom: 12px;" />
                                    <el-skeleton-item variant="text" style="width: 80%; margin-bottom: 8px;" />
                                    <el-skeleton-item variant="text" style="width: 40%;" />
                                </div>
                            </div>
                        </template>

                        <template #default>
                            <template v-if="upcomingContests.length > 0">
                                <div v-for="contest in upcomingContests" :key="contest.id"
                                    class="competition-item upcoming" @click="openContestDialog(contest)">
                                    <div class="competition-content">
                                        <div class="time-badge">
                                            <el-icon>
                                                <Timer />
                                            </el-icon> {{ formatTime(contest.time) }}
                                        </div>
                                        <h3 class="competition-name">{{ contest.title }}</h3>
                                        <div class="competition-description">
                                            {{ contest.text }}
                                        </div>
                                        <div class="competition-link" v-if="contest.link">
                                            <el-link type="primary" :href="contest.link" target="_blank"
                                                @click.stop>比赛链接</el-link>
                                        </div>
                                    </div>
                                </div>
                            </template>
                            <div v-else class="empty-state">
                                <el-empty :image-size="200" description="暂无即将开始的比赛">
                                    <template #description>
                                        <p>暂无即将开始的比赛</p>
                                        <p class="empty-subtitle">敬请期待更多精彩赛事</p>
                                    </template>
                                </el-empty>
                            </div>
                        </template>
                    </el-skeleton>
                </div>
            </div>

            <!-- 已结束的比赛 -->
            <div class="section">
                <h2 class="section-title">
                    <el-icon>
                        <CircleCheck />
                    </el-icon> 已结束
                    <span class="badge">{{ endedContests.length }}</span>
                </h2>

                <div class="competition-list">
                    <template v-if="endedContests.length > 0">
                        <div v-for="contest in endedContests" :key="contest.id" class="competition-item ended"
                            @click="openContestDialog(contest)">
                            <div class="competition-content">
                                <div class="time-badge">
                                    <el-icon>
                                        <Timer />
                                    </el-icon> {{ formatTime(contest.time) }}
                                </div>
                                <div class="ended-overlay">已结束</div>
                                <h3 class="competition-name">{{ contest.title }}</h3>
                                <div class="competition-description">
                                    {{ contest.text }}
                                </div>
                                <div class="competition-link" v-if="contest.link">
                                    <el-link type="info" :href="contest.link" target="_blank" @click.stop>回顾链接</el-link>
                                </div>
                            </div>
                        </div>
                    </template>
                    <div v-else class="empty-state">
                        <el-empty :image-size="200" description="暂无已结束的比赛">
                            <template #description>
                                <p>暂无已结束的比赛</p>
                                <p class="empty-subtitle">历史比赛记录将在这里显示</p>
                            </template>
                        </el-empty>
                    </div>
                </div>
            </div>
        </div>

        <!-- 比赛详情弹窗 -->
        <el-dialog v-model="dialogVisible" :title="selectedContest.title" width="90%" :fullscreen="isMobile">
            <div class="contest-detail">
                <div class="detail-meta">
                    <p><el-icon>
                            <Timer />
                        </el-icon> <strong>比赛时间：</strong> {{ formatTime(selectedContest.time) }}</p>
                    <p v-if="!isContestEnded(selectedContest)"><strong>倒计时：</strong> {{ selectedContest.timeRemaining }}
                    </p>
                    <p><strong>比赛简介：</strong></p>
                    <p>{{ selectedContest.text || '暂无比赛详情描述' }}</p>
                    <div v-if="selectedContest.link" class="detail-link">
                        <el-button type="primary" @click="openLink(selectedContest.link)">
                            {{ isContestEnded(selectedContest) ? '查看回顾' : '前往比赛' }}
                        </el-button>
                    </div>
                </div>
            </div>

            <template #footer>
                <el-button @click="dialogVisible = false">关闭</el-button>
            </template>
        </el-dialog>
    </div>
</template>

<script>
import request from '@/utils/request'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { AlarmClock, Timer, Warning, CircleCheck } from '@element-plus/icons-vue'
import 'dayjs/locale/zh-cn'

// 设置 dayjs 语言为中文
dayjs.locale('zh-cn')

export default {
    name: 'CompetitionPage',
    components: {
        AlarmClock,
        Timer,
        Warning,
        CircleCheck
    },
    data() {
        return {
            contests: [],
            loading: true,
            dialogVisible: false,
            selectedContest: {},
            isMobile: window.innerWidth < 768,
            retryCount: 0,
            refreshInterval: null
        }
    },
    computed: {
        upcomingContests() {
            const now = dayjs();
            return this.contests
                .filter(contest => dayjs(contest.time).isAfter(now))
                .sort((a, b) => dayjs(a.time).diff(dayjs(b.time)));
        },
        endedContests() {
            const now = dayjs();
            return this.contests
                .filter(contest => dayjs(contest.time).isSameOrBefore(now))
                .sort((a, b) => dayjs(b.time).diff(dayjs(a.time)));
        }
    },
    created() {
        this.fetchContests();
        window.addEventListener('resize', this.handleResize);
        // 添加定时刷新
        this.refreshInterval = setInterval(this.fetchContests, 300000); // 每5分钟刷新一次
    },
    beforeUnmount() {
        window.removeEventListener('resize', this.handleResize);
        // 清除定时器
        if (this.refreshInterval) {
            clearInterval(this.refreshInterval);
        }
    },
    methods: {
        openLink(url) {
            if (url) {
                window.open(url, '_blank');
            } else {
                ElMessage({
                    message: '暂无比赛链接',
                    type: 'warning'
                });
            }
        },
        async fetchContests() {
            try {
                this.loading = true;
                const res = await request.get('/contest/list');
                if (res.data && res.data.data) {
                    this.contests = res.data.data;
                    this.retryCount = 0; // 重置重试计数器
                } else {
                    throw new Error('返回数据格式错误');
                }
            } catch (error) {
                console.error('获取比赛列表错误:', error);
                // 如果是网络错误且重试次数小于3次，尝试重试
                if (error.message.includes('Network Error') && this.retryCount < 3) {
                    this.retryCount++;
                    setTimeout(() => {
                        this.fetchContests();
                    }, 3000); // 3秒后重试
                    ElMessage({
                        message: `网络连接失败，正在进行第 ${this.retryCount} 次重试...`,
                        type: 'warning',
                        duration: 2000
                    });
                } else {
                    // 添加更详细的错误信息
                    const errorMessage = error.response?.data?.message
                        || (error.response?.status === 404 ? 'API接口不存在'
                            : error.response?.status === 500 ? '服务器内部错误'
                                : '获取比赛信息失败');

                    ElMessage({
                        message: errorMessage,
                        type: 'error',
                        duration: 5000
                    });
                }
            } finally {
                this.loading = false;
            }
        },
        formatTime(timeStr) {
            if (!timeStr) return '时间待定';
            const time = dayjs(timeStr);
            return time.format('YYYY年MM月DD日 HH:mm');
        },
        formatTimeRemaining(timeStr) {
            const time = dayjs(timeStr);
            const now = dayjs();
            const diff = time.diff(now);

            if (diff < 0) return '已结束';

            const days = Math.floor(diff / (1000 * 60 * 60 * 24));
            const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
            const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

            if (days > 0) return `还有 ${days} 天`;
            if (hours > 0) return `还有 ${hours} 小时`;
            return `还有 ${minutes} 分钟`;
        },
        openContestDialog(contest) {
            this.selectedContest = {
                ...contest,
                timeRemaining: this.formatTimeRemaining(contest.time)
            };
            this.dialogVisible = true;
        },
        isContestEnded(contest) {
            if (!contest.time) return false;
            return dayjs(contest.time).isBefore(dayjs());
        },
        handleResize() {
            this.isMobile = window.innerWidth < 768;
        }
    }
}
</script>

<style scoped>
/* 基础样式 */
.competition-page {
    font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
    color: #333;
    line-height: 1.6;
    padding: 20px;
    max-width: 1200px;
    margin: 0 auto;
}

/* 页面标题 */
.page-header {
    margin-bottom: 40px;
    text-align: center;
    padding: 20px 0;
    /* background: linear-gradient(135deg, #f5f7fa 0%, #e4e8eb 100%); */
    /* border-radius: 12px; */
    /* box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05); */
}

.page-title {
    font-size: 3rem;
    color: var(--dark-color);
    margin-bottom: 2rem;
    text-align: center;
    letter-spacing: 2px;
    font-weight: 700;
    text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
    transition: all 0.3s ease;
}

.page-subtitle {
    font-size: 1.1rem;
    color: #7f8c8d;
    max-width: 600px;
    margin: 0 auto;
    padding: 0 20px;
}

/* 比赛容器 */
.competition-container {
    margin-top: 30px;
}

/* 分区标题 */
.section {
    margin-bottom: 40px;
}

.section-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: #303133;
    display: flex;
    align-items: center;
    margin-bottom: 20px;
    padding-bottom: 10px;
    border-bottom: 2px solid #f0f2f5;
}

.section-title i {
    margin-right: 10px;
    color: #409EFF;
    font-size: 1.3em;
}

.badge {
    background-color: #409EFF;
    color: white;
    border-radius: 12px;
    padding: 4px 10px;
    font-size: 0.9rem;
    margin-left: 10px;
    font-weight: 500;
}

/* 比赛列表 */
.competition-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
    gap: 30px;
    padding: 10px;
}

/* 比赛项 */
.competition-item {
    background-color: #fff;
    border-radius: 16px;
    overflow: hidden;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
    transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
    cursor: pointer;
    border: 1px solid #ebeef5;
    padding: 24px;
    display: flex;
    flex-direction: column;
    gap: 16px;
    position: relative;
}

.competition-item:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.1);
    border-color: #409EFF;
}

.competition-item.upcoming {
    border-left: 4px solid #67C23A;
}

.competition-item.ended {
    border-left: 4px solid #909399;
}

.time-badge {
    background: linear-gradient(135deg, #409EFF, #67C23A);
    color: white;
    padding: 8px 16px;
    border-radius: 20px;
    font-size: 0.9rem;
    display: inline-flex;
    align-items: center;
    gap: 6px;
    margin-bottom: 12px;
    box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2);
}

.ended-overlay {
    background: linear-gradient(135deg, #909399, #606266);
    color: white;
    padding: 8px 16px;
    border-radius: 20px;
    font-size: 0.9rem;
    display: inline-flex;
    align-items: center;
    margin-left: 12px;
    box-shadow: 0 2px 8px rgba(144, 147, 153, 0.2);
}

.competition-content {
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.competition-name {
    font-size: 1.3rem;
    font-weight: 600;
    color: #303133;
    line-height: 1.4;
    margin: 8px 0;
}

.competition-description {
    font-size: 1rem;
    color: #606266;
    line-height: 1.6;
    margin-bottom: 16px;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: 3;
    line-clamp: 3;
    overflow: hidden;
    text-overflow: ellipsis;
}

.competition-link {
    margin-top: auto;
    padding-top: 16px;
    border-top: 1px solid #f0f2f5;
}

.competition-link .el-link {
    font-size: 0.95rem;
    font-weight: 500;
}

/* 空状态 */
.empty-state {
    grid-column: 1 / -1;
    text-align: center;
    padding: 60px 0;
    color: #909399;
    background-color: #fafafa;
    border-radius: 12px;
}

.empty-state i {
    font-size: 3rem;
    margin-bottom: 20px;
    color: #E6A23C;
    opacity: 0.7;
}

.empty-state p {
    font-size: 1.1rem;
    margin-top: 10px;
}

/* 详情弹窗 */
.contest-detail {
    max-height: 70vh;
    overflow-y: auto;
    padding: 24px;
    background: #f8fafc;
    border-radius: 16px;
}

.detail-meta {
    font-size: 1.05rem;
    line-height: 1.8;
    color: #606266;
    background: #fff;
    padding: 24px;
    border-radius: 12px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}

.detail-meta p {
    margin-bottom: 20px;
    display: flex;
    align-items: center;
    gap: 8px;
}

.detail-meta i {
    color: #409EFF;
    font-size: 1.2em;
}

.detail-link {
    margin-top: 32px;
    text-align: center;
}

.detail-link .el-button {
    padding: 12px 24px;
    font-size: 1.05rem;
}

/* 骨架屏 */
.skeleton {
    background-color: #fff;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

/* 移动端适配 */
@media (max-width: 768px) {
    .page-title {
        font-size: 2rem;
    }

    .page-subtitle {
        font-size: 0.95rem;
    }

    .competition-list {
        grid-template-columns: 1fr;
        gap: 20px;
        padding: 5px;
    }

    .section-title {
        font-size: 1.3rem;
    }

    .competition-item {
        padding: 20px;
    }

    .time-badge,
    .ended-overlay {
        padding: 6px 12px;
        font-size: 0.85rem;
    }

    .competition-name {
        font-size: 1.2rem;
    }

    .competition-description {
        font-size: 0.95rem;
    }

    .detail-meta {
        padding: 20px;
    }
}

.empty-subtitle {
    font-size: 0.9rem;
    color: #909399;
    margin-top: 8px;
}
</style>