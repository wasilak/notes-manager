/* global angular */
/* global CodeMirror */

var codeMirror = function ($timeout) {
  return {
    restrict: "E",
    replace: true,
    require: "?ngModel",
    transclude: true,
    scope: {
      syntax: "@",
      theme: "@",
      lineNumbers: "=",
      autoCloseBrackets: "=",
      matchBrackets: "=",
      ngModel: "=",
      breakPoints: "="
    },
    template: '<div class="code-editor"></div>',

    // eslint-disable-next-line no-unused-vars
    link: function (scope, element, attrs, ngModelCtrl, transclude) {
      let editor = CodeMirror(element[0], {
        mode: scope.syntax || "javascript",
        theme: scope.theme || "default",
        autoCloseBrackets: scope.autoCloseBrackets || true,
        matchBrackets: scope.matchBrackets || true,
        lineNumbers: scope.lineNumbers === true ? true : false,
        extraKeys: { "Enter": "newlineAndIndentContinueMarkdownList" },
        lineWrapping: true,
        continueLineComment: true,
        gutters: ["CodeMirror-linenumbers", "breakpoints"]
      });

      scope.breakpoints = [];

      // eslint-disable-next-line no-unused-vars
      scope.$watch('breakPoints', function (current, original) {

        editor.eachLine(function (line) {
          editor.setGutterMarker(line.lineNo(), "breakpoints", null);

        });

        for (let line in current) {
          let info = editor.lineInfo(current[line]);
          editor.setGutterMarker(current[line], "breakpoints", info.gutterMarkers ? null : function () {
            let marker = document.createElement("div");
            marker.style.color = "#822";
            marker.textContent = "●";
            return marker;
          }());
        }
      });

      if (ngModelCtrl) {
        $timeout(function () {
          ngModelCtrl.$render = function () {
            editor.setValue(ngModelCtrl.$viewValue);
          };
        });
      }

      // eslint-disable-next-line no-unused-vars
      transclude(function (clonedEl) {
        //var initialText = clonedEl.text();
        var initialText = scope.ngModel;
        editor.setValue(initialText);

        if (ngModelCtrl) {
          $timeout(function () {
            if (initialText && !ngModelCtrl.$viewValue) {
              ngModelCtrl.$setViewValue(initialText);
            }

            editor.on('change', function () {
              ngModelCtrl.$setViewValue(editor.getValue());


            });
          });
        }
      });

      scope.$on('$destroy', function () {
        editor.off('change');
      });
    }
  };
};

angular.module("app").directive("codeEditor", codeMirror);
