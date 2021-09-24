<template>
  <el-dialog :model-value="show" destroy-on-close width="80%" @close="close">
    <div class="content-center hidden-sm-and-down">
      <div class="content-center">扫描二维码或复制下方字符串到剪贴板并通过剪贴板导入</div>
      <div class="qr-code margin-top">
        <qrcode-vue :size="480" :value="content" level="H"/>
      </div>
      <div class="margin-top">
        <el-input v-model="content" id="vmess-content">
          <template #append>
            <el-button @click="copy" icon="el-icon-document-copy"></el-button>
          </template>
        </el-input>
      </div>
    </div>
    <div class="content-center hidden-md-and-up">
      <div class="content-center">扫描二维码或复制下方字符串到剪贴板并通过剪贴板导入</div>
      <div class="qr-code margin-top">
        <qrcode-vue :size="180" :value="content" level="H"/>
      </div>
      <div class="margin-top">
        <el-input v-model="content" id="vmess-content">
          <template #append>
            <el-button @click="copy" icon="el-icon-document-copy"></el-button>
          </template>
        </el-input>
      </div>
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

    copy() {
      let input = document.querySelector('#vmess-content') as any;
      input.select();

      if (document.execCommand('copy')) {
        this.$message.success("复制成功")
      }
    }
  },
})
</script>

<style lang="scss" scoped>

</style>
