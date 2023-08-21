/* jslint node: true */
"use strict";

angular.module("app").component("cheatsheet",
    {
        controller: function ($scope, $rootScope, $stateParams, ApiService, $state) {
            var vm = this;

            vm.loader = false;


        },
        templateUrl: "/static/app/views/cheatsheet.html"
    }
);
