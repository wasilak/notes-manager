/* jslint node: true */
"use strict";

function ListCtrl($rootScope, $scope, ApiService) {
  var vm = this;

  vm.list = [];

  ApiService.getList().then(function(result) {
    vm.list = result;
  });

  $rootScope.$on('currentNote', function(event, note) {
    vm.note = note;
  });

  vm.search = function() {
    ApiService.getList(vm.listFilter).then(function(result) {
       vm.list = result;
    });
  };
}

ListCtrl.resolve = {
  notes: function($stateParams, ApiService, $rootScope) {
    return ApiService.getList().then(function(result) {
      return result;
    });
  }
};

angular.module("app").controller("ListCtrl", ListCtrl);
