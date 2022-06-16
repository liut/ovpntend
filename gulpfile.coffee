gulp       = require 'gulp'
uglify     = require 'gulp-uglify'
sourcemaps = require 'gulp-sourcemaps'
stylus     = require 'gulp-stylus'
rename     = require 'gulp-rename'
gutil      = require 'gulp-util'
source     = require 'vinyl-source-stream'
buffer     = require 'vinyl-buffer'
browserify = require 'browserify'
del        = require 'del'

paths =
  scripts: [
    './f2e/scripts/home.coffee',
  ]
  stylesheets: [
    './f2e/stylesheets/home.styl',
  ]
  stylus: [
    './node_modules/bootstrap/dist/css',
    './node_modules/bootstrap-select/dist/css',
  ]
  images: [
  ]
  static: './ui/static/'

gulp.task 'build:scripts', gulp.series () ->
  browserify(entries: paths.scripts, debug: true).bundle()
    .pipe(source('bundle.js'))
    .pipe(buffer())
    .pipe(sourcemaps.init(loadMaps: false))
    .pipe(uglify().on('error', gutil.log))
    .pipe(sourcemaps.write())
    .pipe(gulp.dest(paths.static))

gulp.task 'build:stylesheets', gulp.series () ->
  gulp.src(paths.stylesheets)
    .pipe(sourcemaps.init(loadMaps: false))
    .pipe(stylus(compress: true, paths: paths.stylus, 'include css': true))
    .pipe(rename('bundle.css'))
    .pipe(sourcemaps.write())
    .pipe(gulp.dest(paths.static))

gulp.task 'build', gulp.parallel ['build:scripts', 'build:stylesheets']

gulp.task 'watch', gulp.series ['build'], () ->
  gulp.watch paths.scripts, gulp.series ['build:scripts']
  gulp.watch paths.stylesheets, gulp.series ['build:stylesheets']

gulp.task 'clean', gulp.series (cb) ->
  del([paths.static + '*.{js,css,png}'], cb)
