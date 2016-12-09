/* jslint node: true */
"use strict";

function ListCtrl($rootScope, $scope, ApiService) {
  var vm = this;

  vm.list = [];

  ApiService.getList().then(function(result) {
    vm.list = result;
  });
}

ListCtrl.resolve = {
};

angular.module("app").controller("ListCtrl", ListCtrl);
