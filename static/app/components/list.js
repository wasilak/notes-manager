/* jslint node: true */
"use strict";

angular.module("app").component("list", 
  {
    controller: function ListCtrl($rootScope, ApiService) {
      var vm = this;

      vm.list = {
        success: true
      };

      vm.sort = "updated:desc";

      $rootScope.$on('currentNote', function(event, note) {
        vm.note = note;
      });

      vm.search = function() {
        ApiService.getList(vm.listFilter, vm.sort).then(function(result) {
          vm.list = result;
        });
      };

      vm.setSort = function() {
        vm.search();
      };

      vm.search();
    },
    templateUrl: "/static/app/views/list.html"
});
