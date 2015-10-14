var app = angular.module("AdminApp", ["ngResource", "ui.bootstrap", "countTo"]);

app.config(function ($provide) {
    var progressDecorator = function ($delegate) {
        var directive = $delegate[0];
        var compile = directive.compile;
        var link = directive.link;
        directive.compile = function () {
            compile.apply(this, arguments);
            return function (scope, elem, attr, ctrl) {
                link.apply(this, arguments);
                if (angular.isDefined(attr.dynamicMax)) {
                    attr.$observe('dynamicMax', function (max) {
                        scope.max = max;
                        scope.percent = +(100 * scope.value / max).toFixed(2);
                    });
                }
            };
        };
        return $delegate;
    };
    $provide.decorator('progressbarDirective', progressDecorator);
    $provide.decorator('barDirective', progressDecorator);
});

app.factory("Repositories", function ($resource) {
    return $resource("/repo");
});

app.factory("Repository", function ($resource) {
    return $resource('/repo', {}, {
        delete: {method: "DELETE", data: {url: "@url"}},
        save: {method: "POST", data: {type: "@type", url: "@url"}}
    });
});

app.controller("AllRepo", function ($scope, Repositories, Repository) {
    Repositories.query(function (data) {
        $scope.repositories = data;
        $scope.$parent.$broadcast("InitProgress", {"countMax": $scope.repositories.length});
    });
    $scope.deleteRepo = function (repo) {
        if (confirm("Are you really want delete repository?")) {
            Repository.delete({url: repo.url});
            var index = -1;
            var comArr = eval($scope.repositories);
            for (var i = 0; i < comArr.length; i++) {
                if (comArr[i].url === repo.url) {
                    index = i;
                    break;
                }
            }
            if (index === -1) {
                alert("Something gone wrong");
            }
            $scope.repositories.splice(index, 1);
        }
    }
});

app.controller("AddRepoModal", function ($scope, $modal, $log, Repository) {
    $scope.animationsEnabled = true;
    $scope.repository = Repository;
    $scope.open = function (event) {
        event.preventDefault();
        var modalInstance = $modal.open({
            animation: $scope.animationsEnabled,
            templateUrl: 'addRepoContent.html',
            controller: 'AddRepoModalInstance',
            resolve: {
                repository: function () {
                    return $scope.repository;
                }
            }
        });
        modalInstance.result.then(function (data) {
            Repository.save({type: data.type, url: data.url});
        });
    };
});

app.controller("AddRepoModalInstance", function ($scope, $modalInstance) {
    $scope.submitForm = function () {
        $modalInstance.close($scope.repository);
    };
    $scope.cancel = function () {
        $modalInstance.dismiss('cancel');
    };
});

app.factory("Packages", function ($resource) {
    return $resource("/package/update", {});
});

app.controller("PackageUpdate", function ($scope, Packages) {
    $scope.update = function (event) {
        event.preventDefault();
        $scope.$parent.$broadcast("ClearProgress");
        var source = new EventSource("/package/update");
        source.onmessage = function (event) {
            $scope.$apply(function () {
                $scope.$parent.$broadcast("RunProgress", {"current": event.data});
            });
        };
        source.onopen = function (event) {
            console.log("eventsource: open");
        };
        source.onerror = function (event) {
            if (event.target.readyState === 0) {
                console.log("eventsource: close");
                source.close()
            } else {
                console.log("eventsource: error");
            }
        };
    };
});

app.controller("ProgressBar", function ($scope) {
    $scope.$parent.$on("InitProgress", function (event, data) {
        $scope.countMax = data.countMax;
    });
    $scope.$parent.$on("RunProgress", function (event, data) {
        $scope.countFrom = 0;
        $scope.countTo = data.current;
        $scope.progressValue = data.current;
    });
    $scope.$parent.$on("ClearProgress", function (event, data) {
        $scope.countFrom = 0;
        $scope.countTo = 0;
        $scope.progressValue = 0;
    });
});