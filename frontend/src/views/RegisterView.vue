<template>
  <div>
    <img src="../assets/bg.jpg" style="float: left;display: inline-block;width: 50%;margin-top: 5%">
    <div class="register-container">
      <h2>注册页面</h2>
      <!-- <button @click="toRegister" style="background-color: grey;">点我注册</button> -->
      <form @submit.prevent="register">
        <div class="form-group">
          <label for="name">用户名:</label>
          <input v-model="name" type="text" id="name" placeholder="用户名" required>
        </div>
        <div>
          <label for="address">钱包地址:</label>
          <input id="address" type="text" v-model="address"/>
        </div>
        <div class="form-group">
          <label for="password">密码:</label>
          <input v-model="password1" type="password" id="password" placeholder="密码" required>
        </div>
        <div class="form-group">
          <label for="password">确认密码:</label>
          <input v-model="password2" type="password" id="password" placeholder="确认密码" required>

        </div>
        <button type="submit" @click="this.$router.push({path: '/'})">注册</button>

      </form>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      name: '',
      address: '',
      password1: '',
      password2: '',
    };
  },
  methods: {
    register() {
      // 处理登录逻辑
      if (this.address && (this.password1 === this.password2)) {
        let formData = new FormData();
        formData.append("username", this.name);
        formData.append("address", this.address);
        formData.append("password", this.password1);
        axios.post("/user/register", formData, {
          headers: {'Content-Type': 'application/x-www-form-urlencoded'},
        }).then(function (resp) {
          switch (resp.status) {
            case 200:
              this.error = "注册成功，请登录";
              break;
            case 600:
              let obj = JSON.parse(resp.data);
              this.error = obj.desc;
          }
        }).catch(function (err) {
          console.log(err)
        })
      }
    },

  }
};
</script>

<style scoped>
.register-container {
  max-width: 400px;
  padding: 20px;
  background-color: #f5f5f5;
  margin: 5%;
  float: right;
}

.form-group {
  margin-bottom: 16px;
}

label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
  float: left;
}

input[type="text"],
input[type="password"] {
  width: 100%;
  padding: 8px;
  font-size: 16px;
  border-radius: 4px;
  border: 1px solid #ccc;
}

button[type="submit"] {
  display: block;
  width: 100%;
  padding: 12px;
  font-size: 16px;
  font-weight: bold;
  color: #fff;
  background-color: #007bff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
</style>
  