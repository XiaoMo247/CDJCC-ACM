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
                    <i class="el-icon-alarm-clock"></i> 即将开始
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
                                    <div class="competition-image"
                                        :style="{ backgroundImage: `url(${getUpcomingImage(contest.id)})` }">
                                        <div class="time-badge">
                                            <i class="el-icon-time"></i> {{ formatTime(contest.time) }}
                                        </div>
                                    </div>
                                    <div class="competition-content">
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
                                <i class="el-icon-warning-outline"></i>
                                <p>暂无即将开始的比赛</p>
                            </div>
                        </template>
                    </el-skeleton>
                </div>
            </div>

            <!-- 已结束的比赛 -->
            <div class="section">
                <h2 class="section-title">
                    <i class="el-icon-finished"></i> 已结束
                    <span class="badge">{{ endedContests.length }}</span>
                </h2>

                <div class="competition-list">
                    <template v-if="endedContests.length > 0">
                        <div v-for="contest in endedContests" :key="contest.id" class="competition-item ended"
                            @click="openContestDialog(contest)">
                            <div class="competition-image"
                                :style="{ backgroundImage: `url(${getEndedImage(contest.id)})` }">
                                <div class="ended-overlay">已结束</div>
                                <div class="time-badge">
                                    <i class="el-icon-time"></i> {{ formatTime(contest.time) }}
                                </div>
                            </div>
                            <div class="competition-content">
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
                        <i class="el-icon-warning-outline"></i>
                        <p>暂无已结束的比赛</p>
                    </div>
                </div>
            </div>
        </div>

        <!-- 比赛详情弹窗 -->
        <el-dialog v-model="dialogVisible" :title="selectedContest.title" width="90%" :fullscreen="isMobile">
            <div class="contest-detail">
                <div class="detail-image"
                    :style="{ backgroundImage: `url(${isContestEnded(selectedContest) ? getEndedImage(selectedContest.id) : getUpcomingImage(selectedContest.id)})` }">
                </div>

                <div class="detail-meta">
                    <p><i class="el-icon-time"></i> <strong>比赛时间：</strong> {{ formatTime(selectedContest.time) }}</p>
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

export default {
    name: 'CompetitionPage',
    data() {
        return {
            contests: [],
            loading: true,
            dialogVisible: false,
            selectedContest: {},
            isMobile: window.innerWidth < 768,
            upcomingImages: [
                'https://img.freepik.com/free-vector/gradient-hackathon-logo_23-2149324865.jpg',
                'https://img.freepik.com/free-vector/hand-drawn-science-education-background_23-2148499325.jpg',
                'https://img.freepik.com/free-vector/gradient-coding-company-logo_23-2148821801.jpg'
            ],
            endedImages: [
                'https://img.freepik.com/free-vector/gradient-hackathon-background_23-2149324866.jpg',
                'https://img.freepik.com/free-vector/gradient-coding-background_23-2149159430.jpg',
                'https://img.freepik.com/free-vector/gradient-coding-logo-template_23-2148807059.jpg'
            ]
        }
    },
    computed: {
        upcomingContests() {
            const now = new Date();
            return this.contests.filter(contest => {
                return new Date(contest.time) > now;
            }).sort((a, b) => new Date(a.time) - new Date(b.time));
        },
        endedContests() {
            const now = new Date();
            return this.contests.filter(contest => {
                return new Date(contest.time) <= now;
            }).sort((a, b) => new Date(b.time) - new Date(a.time));
        }
    },
    created() {
        this.fetchContests();
        window.addEventListener('resize', this.handleResize);
    },
    beforeUnmount() {
        window.removeEventListener('resize', this.handleResize);
    },
    methods: {
        openLink(url) {
            if (url) {
                window.open(url, '_blank');
            } else {
                this.$message.warning('暂无比赛链接');
            }
        },
        async fetchContests() {
            try {
                this.loading = true;
                const res = await request.get('/contest/list');
                this.contests = res.data.data || [];
            } catch (error) {
                this.$message.error('获取比赛信息失败');
            } finally {
                this.loading = false;
            }
        },
        formatTime(timeStr) {
            if (!timeStr) return '时间待定';
            return dayjs(timeStr).format('YYYY年MM月DD日 HH:mm');
        },
        openContestDialog(contest) {
            this.selectedContest = contest;
            this.dialogVisible = true;
        },
        isContestEnded(contest) {
            if (!contest.time) return false;
            return new Date(contest.time) <= new Date();
        },
        getUpcomingImage(id) {
            // 根据ID选择不同的即将开始比赛图片
            return this.upcomingImages[id % this.upcomingImages.length];
        },
        getEndedImage(id) {
            // 根据ID选择不同的已结束比赛图片
            return this.endedImages[id % this.endedImages.length];
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
    font-weight: bold;
    color: #2c3e50;
    margin-bottom: 10px;
    text-shadow: 1px 1px 2px rgba(0,0,0,0.1);
    background: linear-gradient(90deg, #409EFF, #67C23A);
    -webkit-background-clip: text;
    background-clip: text;
    color: transparent;
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
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 25px;
}

/* 比赛项 */
.competition-item {
    background-color: #fff;
    border-radius: 12px;
    overflow: hidden;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
    cursor: pointer;
    border: 1px solid #ebeef5;
}

.competition-item:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 25px rgba(0, 0, 0, 0.12);
}

.competition-item.upcoming .competition-image {
    border-bottom: 4px solid #67C23A;
}

.competition-item.ended .competition-image {
    border-bottom: 4px solid #909399;
    filter: grayscale(20%);
}

.competition-image {
    height: 160px;
    background-size: cover;
    background-position: center;
    position: relative;
    transition: all 0.3s;
}

.competition-item:hover .competition-image {
    transform: scale(1.02);
}

.time-badge {
    position: absolute;
    bottom: 10px;
    left: 10px;
    background-color: rgba(0, 0, 0, 0.7);
    color: white;
    padding: 6px 12px;
    border-radius: 20px;
    font-size: 0.8rem;
    backdrop-filter: blur(2px);
}

.ended-overlay {
    position: absolute;
    top: 10px;
    right: 10px;
    background-color: rgba(0, 0, 0, 0.7);
    color: white;
    padding: 6px 12px;
    border-radius: 20px;
    font-size: 0.8rem;
    backdrop-filter: blur(2px);
}

.competition-content {
    padding: 20px;
}

.competition-name {
    font-size: 1.2rem;
    font-weight: 600;
    margin-bottom: 12px;
    color: #303133;
    line-height: 1.4;
}

.competition-description {
    font-size: 0.95rem;
    color: #606266;
    margin-bottom: 15px;
    display: -webkit-box;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.6;
}

.competition-link {
    margin-top: 15px;
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
}

.detail-image {
    width: 100%;
    height: 250px;
    border-radius: 8px;
    background-size: cover;
    background-position: center;
    margin-bottom: 25px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.detail-meta {
    font-size: 1rem;
    line-height: 1.8;
    color: #606266;
    padding: 0 10px;
}

.detail-meta p {
    margin-bottom: 18px;
}

.detail-meta i {
    margin-right: 10px;
    color: #409EFF;
    font-size: 1.1em;
}

.detail-link {
    margin-top: 30px;
    text-align: center;
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
    }

    .section-title {
        font-size: 1.3rem;
    }

    .detail-image {
        height: 180px;
    }
    
    .competition-image {
        height: 140px;
    }
    
    .empty-state {
        padding: 40px 20px;
    }
    
    .empty-state i {
        font-size: 2.5rem;
    }
}
</style>