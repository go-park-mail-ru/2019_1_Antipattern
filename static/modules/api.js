'use strict';

(function() {
	const noop = () => null;
  const ajax = window.AjaxModule;
  const apiPreffix = "";
	class API {
    _sendRequest(url, body, callback) {
      ajax.doPost( {
       callback: (xhr) => {
          try {
            const object = JSON.parse(xhr.response);
          }
          catch (SyntaxError) {
            if (callback) {
              callback('broken_response', null);
            }
            
          }
          if (!callback) {
            return
          }
          if (object.type != 'reg') {
            callback('wrong_response', null)
            return
          }
          callback(object.status, object.payload);
        }, 
        path : apiPreffix + url, 
        body : body 
      });
    }

		register(login, email, password, name, callback) {
      const url = '/api/register';
      const data = {
        login : login,
        email : email,
        password : password,
        name : name
      };
      this._sendRequest(url, data, callback)
    }

    authorize(login, password, callback) {
      url = '/api/auth';
      const data = {
        login : login,
        password : password
      };
      this._sendRequest(url, data, callback)
    }

    updateProfile(login, password, callback) {
      url = '/api/update';
    }
    
    getUserInfo(login, password, callback) {
      url = '/api/get_info'
    
    }
    
    isAuthorized(callback) {
      url = '/api/is_authorized'

    }
	}
	window.API = new API();
})();

