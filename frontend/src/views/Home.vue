<template>
  <div class="content-center el-container">
    <div class="el-main el-main-md-and-up hidden-sm-and-down">
      <template v-for="item in data" v-bind:key="item.id">
        <div class="endpoint-box el-row margin-bottom-2x" v-on:mouseover="item.show_delete_button = true"
             v-on:mouseout="item.show_delete_button = false">
          <div class="el-col-2">{{ getCloud(item.cloud) }}</div>
          <div class="el-col-2">{{ getEndpoint(item.endpoint) }}</div>
          <div class="el-col-4">{{ item.host }}</div>
          <div class="el-col-2">{{ item.rate ? item.rate : '-' }}</div>
          <div class="el-col-2">{{ getMode(item.transport_type) }}</div>
          <div class="el-col-xl-10 el-col-lg-10 el-col-md-10">{{ item.remark ? item.remark : '-' }}</div>
          <div class="el-col-xl-2 el-col-lg-2 el-col-md-2">
            <el-button type="primary" :loading="item.downloading" @click="download(item)">生成配置</el-button>
            <el-button type="danger" class="delete-button" icon="el-icon-delete" circle
                       @click="removeItem = item; showRemoveModal = true" v-show="item.show_delete_button"></el-button>
          </div>
        </div>
      </template>

      <div class="endpoint-box el-row margin-bottom-2x" v-if="0 === data.length">
        <div class="el-col-24">暂时没有可用的节点列表</div>
      </div>

      <div class="el-row">
        <div class="el-col-24">
          <el-button type="success" size="medium" @click="showNewModal = true">添加节点</el-button>
          <el-button type="primary" size="medium" @click="showDownloadModal = true">下载客户端</el-button>
        </div>
      </div>
    </div>

    <div class="el-main el-main-sm-and-down hidden-md-and-up">
      <template v-for="item in data" v-bind:key="item.id">
        <div class="endpoint-box el-row margin-bottom" v-show="!item.show_generate_menu">
          <div class="el-col-12">{{ item.host }}</div>
          <div class="el-col-9">{{ getMode(item.transport_type) }}</div>
          <div class="el-col-3">
            <el-button type="primary" icon="el-icon-menu" circle @click="item.show_generate_menu = true"></el-button>
          </div>
        </div>
        <div class="endpoint-box el-row margin-bottom" v-show="item.show_generate_menu">
          <div class="el-col-21">
            <el-button type="primary" :loading="item.downloading" @click="download(item)">生成配置</el-button>
            <el-button type="danger" @click="removeItem = item; showRemoveModal = true">删除</el-button>
          </div>
          <div class="el-col-3">
            <el-button type="primary" icon="el-icon-menu" circle @click="item.show_generate_menu = false"></el-button>
          </div>
        </div>
      </template>

      <div class="endpoint-box el-row margin-bottom" v-if="0 === data.length">
        <div class="el-col-24">暂时没有可用的节点列表</div>
      </div>

      <div class="el-row">
        <div class="el-col-24">
          <el-button type="success" size="medium" @click="showNewModal = true">添加节点</el-button>
          <el-button type="primary" size="medium" @click="showDownloadModal = true">下载客户端</el-button>
        </div>
      </div>
    </div>
  </div>

  <Remove v-model:show="showRemoveModal" v-bind:loading="removing" v-on:confirm="remove"></Remove>

  <NewV2rayEndpoint v-model:show="showNewModal" v-on:success="load"></NewV2rayEndpoint>

  <Download v-model:show="showDownloadModal"></Download>

  <QRCode v-model:show="showQRCodeModal" v-bind:content="v2rayNgVMessString"></QRCode>
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
import NewV2rayEndpoint from "@/components/NewV2rayEndpoint.vue"
import Download from "@/components/Download.vue"
import QRCode from "@/components/QRCode.vue"

const md5 = require('md5')

export default defineComponent({
  name: 'Home',
  components: {QRCode, Download, NewV2rayEndpoint, Remove},
  data() {
    return {
      loading: true,
      data: new Array<V2rayEndpointListItem>(),
      showRemoveModal: false,
      removeItem: new V2rayEndpointListItem(),
      removing: false,
      showNewModal: false,
      showDownloadModal: false,
      showQRCodeModal: false,
      v2rayNgVMessString: "",
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
      let map = {1: "TCP", 2: "WebSocket", 3: "KCP", 4: "HTTP2"} as any

      return map[transportType]
    },
    download(item: V2rayEndpointListItem) {
      item.downloading = true

      let form = new V2rayEndpointDownloadForm()
      form.id = item.id

      axios.post(API_V2RAY_ENDPOINT_DOWNLOAD, form).then((response: AxiosResponse<V2rayEndpointDownloadResponse>) => {
        item.downloading = false

        if (0 != response.data.code) {
          this.$message.error(response.data.message)

          return
        }

        this.v2rayNgVMessString = response.data.data.content
        this.showQRCodeModal = true
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

.el-main.el-main-md-and-up {
  --el-main-padding: 30px;

  .endpoint-box {
    background-color: #fff;
    padding: 25px;

    .delete-button {
      position: absolute;
      top: -15px;
      right: -15px;
    }
  }
}

.el-main.el-main-sm-and-down {
  --el-main-padding: 15px;

  .endpoint-box {
    background-color: #fff;
    padding: 15px;
  }
}
</style>
