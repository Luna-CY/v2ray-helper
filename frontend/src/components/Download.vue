<template>
  <el-dialog width="60%" :model-value="show" @close="close" destroy-on-close>
    <el-table :data="data" border style="width: 100%">
      <el-table-column prop="os" label="系统类型"></el-table-column>
      <el-table-column prop="client" label="客户端"></el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="scope">
          <el-link type="primary" :underline="false" :href="scope.row.git" target="_blank">GitHub</el-link>
          <el-link type="primary" :underline="false" :href="scope.row.address" class="margin-left" target="_blank">下载
          </el-link>
          <el-link type="primary" :underline="false" :href="scope.row.helper" class="margin-left">配置帮助</el-link>
        </template>
      </el-table-column>
    </el-table>
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

  data() {
    return {
      data: [
        {
          os: "Windows",
          client: "V2rayN",
          git: "https://github.com/2dust/v2rayN/releases",
          address: "https://github.com/2dust/v2rayN/releases/download/4.20/v2rayN-Core.zip",
          helper: "#/helper/v2ray-n",
        },
        {
          os: "MacOSX",
          client: "V2rayX",
          git: "https://github.com/Cenmrev/V2RayX/releases",
          address: "https://github.com/Cenmrev/V2RayX/releases/download/v1.5.1/V2RayX.app.zip",
          helper: "#/helper/v2ray-x",
        },
        {
          os: "Android",
          client: "V2rayNG",
          git: "https://github.com/2dust/v2rayNG/releases",
          address: "https://github.com/2dust/v2rayNG/releases/download/1.6.18/v2rayNG_1.6.18.apk",
          helper: "#/helper/v2ray-ng",
        },
      ],
    }
  },

  methods: {
    close() {
      this.$emit('update:show', false)
    },
  },
})
</script>

<style scoped lang="scss">
.download-client-box-md-and-up {
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

.download-client-box-sm-and-down {
  display: flex;
  flex-flow: column;
  align-items: center;

  .download-client-item {
    width: 100%;
    display: flex;
    flex-flow: column;
    align-items: center;

    .client-image {
      width: 64px;
      height: 64px;
    }
  }
}
</style>
