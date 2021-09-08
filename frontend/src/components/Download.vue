<template>
  <el-dialog width="60%" :model-value="show" @close="close" destroy-on-close>
    <div class="download-client-box">
      <div class="download-client-item">
        <a href="https://github.com/2dust/v2rayN/releases/download/4.20/v2rayN-Core.zip">
          <el-image class="client-image" src="/img/client/v2rayN.ico"></el-image>
        </a>
        <div>V2rayN</div>
      </div>
      <div class="download-client-item">
        <a href="https://github.com/Cenmrev/V2RayX/releases/download/v1.5.1/V2RayX.app.zip">
          <el-image class="client-image" src="/img/client/vx128.png"></el-image>
        </a>
        <div>V2rayX</div>
      </div>
      <div class="download-client-item">
        <a href="https://github.com/2dust/v2rayNG/releases/download/1.6.18/v2rayNG_1.6.18.apk">
          <el-image class="client-image" src="/img/client/v2rayNG.png"></el-image>
        </a>
        <div>V2rayNG</div>
      </div>
    </div>
  </el-dialog>
</template>

<script lang="ts">
import {defineComponent} from "vue"

export default defineComponent({
  name: "Download",

  props: {
    show: Boolean
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
.download-client-box {
  height: 300px;
  width: 100%;
  display: flex;
  justify-content: space-around;
  align-items: center;

  .download-client-item {
    width: 100%;
    display: flex;
    flex-flow: column;
    align-items: center;

    .client-image {
      width: 128px;
      height: 128px;
    }
  }
}
</style>