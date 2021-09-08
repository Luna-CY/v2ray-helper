<template>
  <el-dialog width="60%" :model-value="show" @close="close" destroy-on-close>
    <div class="qr-code-box">
      <div class="content-center">扫描二维码或复制下方字符串到剪贴板并通过剪贴板导入</div>
      <div class="qr-code margin-top">
        <qrcode-vue :value="content" :size="350" level="H"/>
      </div>
      <div class="margin-top">{{ content }}</div>
    </div>
  </el-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue"
import QrcodeVue from 'qrcode.vue'

export default defineComponent({
  name: "QRCode",

  components: {QrcodeVue},

  props: {
    show: Boolean,
    content: String,
  },

  emits: ['update:show'],

  methods: {
    close() {
      this.$emit('update:show', false)
    },
  },
})
</script>

<style scoped lang="scss">
.qr-code {
  text-align: center;
}
</style>