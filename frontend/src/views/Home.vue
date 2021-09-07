<template>
  <div class="home el-container">
    <div class="el-main">
      <template v-for="item in data" v-bind:key="item.id">
        <div class="endpoint-box el-row">
          <div class="el-col-lg-2 el-col-md-1 el-col-sm-2">{{ getCloud(item.cloud) }}</div>
          <div class="el-col-lg-2 el-col-md-1 el-col-sm-2">{{ getEndpoint(item.endpoint) }}</div>
          <div class="el-col-lg-4 el-col-md-4 el-col-sm-4">{{ item.host }}</div>
          <div class="el-col-lg-2 el-col-md-2 el-col-sm-0">{{ item.rate ? item.rate : '-' }}</div>
          <div class="el-col-lg-2 el-col-md-2 el-col-sm-0">{{ getMode(item.transport_type) }}</div>
          <div class="el-col-lg-8 el-col-md-2 el-col-sm-2">{{ item.remark ? item.remark : '-' }}</div>
          <div class="el-col-lg-4 el-col-md-2 el-col-sm-2">
            <el-button type="primary" :loading="downloading" @click="download(item, 1)">V2rayX</el-button>
            <el-button type="primary" :loading="downloading" class="margin-left" @click="download(item, 2)">V2rayNG
            </el-button>
            <el-button type="danger" class="margin-left" @click="removeItem = item; showRemoveModal = true">删除
            </el-button>
          </div>
        </div>
      </template>

      <div class="endpoint-box el-row" v-if="0 === data.length">
        <div class="el-col-lg-24 el-col-md-24 el-col-sm-24">暂时没有可用的节点列表</div>
      </div>

      <div class="el-row">
        <div class="el-col-lg-24 el-col-md-24 el-col-sm-24">
          <el-button type="success" size="medium" @click="showNewModal = true">添加节点</el-button>
        </div>
      </div>
    </div>
  </div>

  <Remove v-model:show="showRemoveModal" v-bind:loading="removing" v-on:confirm="remove"></Remove>

  <NewV2rayEndpoint v-model:show="showNewModal" v-on:success="load"></NewV2rayEndpoint>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import axios, {AxiosResponse} from "axios"
import {API_V2RAY_ENDPOINT_LIST, V2rayEndpointListItem, V2rayEndpointListResponse} from "@/api/v2ray_endpoint_list"
import {
  API_V2RAY_ENDPOINT_DOWNLOAD,
  V2rayEndpointDownloadForm,
  V2rayEndpointDownloadResponse
} from "@/api/v2ray_endpoint_download"
import {API_V2RAY_ENDPOINT_REMOVE, V2rayEndpointRemoveForm} from "@/api/v2ray_endpoint_remove"
import {BaseResponse} from "@/api/base"
import Remove from "@/components/Remove.vue"
import NewV2rayEndpoint from "@/components/NewV2rayEndpoint.vue";

const md5 = require('md5')

export default defineComponent({
  name: 'Home',
  components: {NewV2rayEndpoint, Remove},
  data() {
    return {
      loading: true,
      data: new Array<V2rayEndpointListItem>(),
      downloading: false,
      showRemoveModal: false,
      removeItem: new V2rayEndpointListItem(),
      removing: false,
      showNewModal: false,
    }
  },

  mounted() {
    this.load()
  },

  methods: {
    load() {
      this.loading = true

      axios.get(API_V2RAY_ENDPOINT_LIST).then((response: AxiosResponse<V2rayEndpointListResponse>) => {
        this.loading = false

        if (0 != response.data.code) {
          this.$message.error(response.data.message)

          return
        }

        this.data = response.data.data.data
      })
    },
    getCloud(cloud: number) {
      let map = {1: "Vultr", 2: "阿里云", 3: "腾讯云", 4: "华为云"} as any

      return map[cloud]
    },

    getEndpoint(endpoint: number) {
      let map = {1: "日本", 2: "香港"} as any

      return map[endpoint]
    },
    getMode(transportType: number) {
      let map = {1: "TCP", 2: "WebSocket"} as any

      return map[transportType]
    },
    download(item: V2rayEndpointListItem, type: number) {
      this.downloading = true

      let form = new V2rayEndpointDownloadForm()
      form.id = item.id
      form.type = type

      axios.post(API_V2RAY_ENDPOINT_DOWNLOAD, form).then((response: AxiosResponse<V2rayEndpointDownloadResponse>) => {
        this.downloading = false

        if (0 != response.data.code) {
          this.$message.error(response.data.message)

          return
        }

        console.log(response.data.data.content)

        let element = document.createElement('a');
        element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(response.data.data.content));
        element.setAttribute('download', "test-v2rayX-config.json");

        element.style.display = 'none';
        document.body.appendChild(element);

        element.click();

        document.body.removeChild(element);
      })
    },
    remove(password: string) {
      const body = new V2rayEndpointRemoveForm()
      body.id = this.removeItem.id
      body.password = md5(password)

      this.removing = true
      axios.post(API_V2RAY_ENDPOINT_REMOVE, body).then((response: AxiosResponse<BaseResponse>) => {
        this.removing = false

        if (0 != response.data.code) {
          this.$message.error(response.data.message)

          return
        }

        this.showRemoveModal = false
        this.$message.success("删除成功")

        this.load()
      })
    },
  },
});
</script>

<style lang="scss">
body {
  background-color: #2d3a4b;
}
</style>

<style lang="scss">
.el-row {
  align-items: center;
}

.home {
  width: 80%;
  margin: auto;
  text-align: center;
}

.endpoint-box {
  background-color: #fff;
  margin: 30px;
  padding: 25px;
}
</style>
