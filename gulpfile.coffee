gulp       = require 'gulp'
uglify     = require 'gulp-uglify'
sourcemaps = require 'gulp-sourcemaps'
stylus     = require 'gulp-stylus'
rename     = require 'gulp-rename'
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
  dest: './ui/static/'

gulp.task 'build', ['build:scripts', 'build:stylesheets', 'build:images']

gulp.task 'build:scripts', () ->
  browserify(entries: paths.scripts, debug: true).bundle()
    .pipe(source('bundle.js'))
    .pipe(buffer())
    .pipe(sourcemaps.init(loadMaps: false))
    .pipe(uglify())
    .pipe(sourcemaps.write())
    .pipe(gulp.dest(paths.dest))

gulp.task 'build:stylesheets', () ->
  gulp.src(paths.stylesheets)
    .pipe(sourcemaps.init(loadMaps: false))
    .pipe(stylus(compress: true, paths: paths.stylus, 'include css': true))
    .pipe(rename('bundle.css'))
    .pipe(sourcemaps.write())
    .pipe(gulp.dest(paths.dest))

gulp.task 'build:images', () ->
  gulp.src(paths.images)
    .pipe(gulp.dest(paths.dest))

gulp.task 'watch', ['build'], () ->
  gulp.watch(paths.scripts, ['build:scripts'])
  gulp.watch(paths.stylesheets, ['build:stylesheets'])

gulp.task 'clean', (cb) ->
  del([paths.dest + '*.{js,css,png}'], cb)
