import axios from 'axios';

const API_URL = 'http://localhost:8080/api/';

class AuthService {
  login(user: any) {
    return axios
      .post(API_URL + 'login', {
        username: user.username,
        password: user.password
      })
      .then(response => {
        if (response.data.token) {
          localStorage.setItem('token', response.data.token);
        }
        return response.data;
      });
  }

  logout() {
    localStorage.removeItem('token');
  }

  register(user: any) {
    return axios.post(API_URL + 'register', {
      username: user.username,
      password: user.password
    });
  }

  getToken() {
    return localStorage.getItem('token');
  }
}

export default new AuthService();
