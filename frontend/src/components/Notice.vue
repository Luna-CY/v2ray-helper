<template>
  <el-dialog :model-value="show" destroy-on-close width="80%" @close="close">
    <template v-for="notice in $store.getters.local.notice_list" v-bind:key="notice.id">
      <el-alert :title="notice.title" :type="{1: 'info', 2: 'warning', 3: 'error'}[notice.type]"
                show-icon :closable="false" class="margin-bottom">
        {{ (new Date(notice.time * 1000)).format('yyyy-MM-dd hh:mm:ss') }} - {{ notice.message }}
      </el-alert>
    </template>
    <div class="content-center margin-bottom" v-if="0 === $store.getters.local.notice_list.length">没有通知信息</div>
    <div class="content-center" v-if="0 !== $store.getters.local.notice_list.length">
      <el-button type="primary" @click="clean">清空通知</el-button>
    </div>
  </el-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue"
import axios, {AxiosResponse} from "axios"
import {API_CLEAN_NOTICE} from "@/api/clean_notice"
import {BaseResponse} from "@/api/base"
import {NoticeListItem} from "@/api/meta_info"

export default defineComponent({
  name: "Notice",

  props: {
    show: Boolean,
  },

  emits: ['update:show'],

  methods: {
    close() {
      this.$emit('update:show', false)
    },

    clean() {
      axios.post(API_CLEAN_NOTICE).then((response: AxiosResponse<BaseResponse>) => {
        if (0 != response.data.code) {
          this.$message.error(response.data.message)

          return
        }

        let story = this.$store.getters.local
        story.notice_list = new Array<NoticeListItem>()

        this.$store.commit('local', story)
      })
    }
  }
})
</script>

<style scoped lang="scss">

</style>
