var gorss = angular.module('gorss', ["feedService"]);

gorss.config(function($locationProvider, $routeProvider) {
  $routeProvider.
  when('/', {
    templateUrl: 'templates/home.html',
    controller: "HomeCtrl"
  }).
  otherwise({redirectTo: '/404'});
});

var feedService = angular.module("feedService", ["ngResource"]);

feedService.factory('Feed', ['$resource',
  function($resource){
    return $resource('api/feed', {}, {
      query: {method:'GET', isArray:true}
    });
  }]);

gorss.controller("HomeCtrl", ["$scope", "Feed", "$http", "$window", function($scope, Feed, $http, $window) {
  $scope.loading = true;
  $scope.feeds = Feed.query(function() {
    $scope.loading = false;
  }, function() {
    $scope.loading = false;
  });

  $scope.update = function() {
    $scope.loading = true;
    $scope.feeds = Feed.query(function(data){
      $scope.loading = false;
    });
  };

  $scope.addFeed = function() {
    $http({
      method: "POST",
      url: "/api/feed",
      data: {url: $scope.url}
    }).success(function(data) {
      $scope.feeds.push(data);
    }).error(function(data) {
      $(".alert-danger").fadeIn(500);
      setTimeout(function() {
        $scope.$apply(function(){
          $(".alert-danger").fadeOut(500);
        });
      }, 2000);
    });

    $scope.url = "";
  }

  $scope.delete = function(feedIndex) {
    var feed = $scope.feeds[feedIndex];
    
    $http({
      method: "DELETE",
      url: "/api/feed/",
      data: {url: feed.UpdateURL}
    })

    $scope.feeds.splice(feedIndex, 1);
  }

  $scope.markUnread = function(link, feedIndex, itemIndex) {
    var item = $scope.feeds[feedIndex].Items[itemIndex];

    $http({
      method: "POST",
      url: "/api/feed/read",
      data: {id: item.ID}
    });

    item.Read = true

    $window.open(link);
  };

}]);