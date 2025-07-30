const out: {
  catchName: string
  time: Date
  name: string
  message: string
  stack: object
}[] = []

function writeError(
  catchName: string,
  error: { name: string; message: string; stack: object },
) {
  if (!error) return
  out.unshift({
    catchName,
    time: new Date(),
    name: error.name,
    message: error.message,
    stack: error.stack,
  })
}

function $reportErrorGuard(app, router) {
  // app.config.errorHandler = function (err, instance, info) {
  //   reportError({
  //     error: info,
  //     message: err.message,
  //     instance,
  //     filename: err.stack,
  //   })

  //   writeError('App errorHandler', {
  //     error: info,
  //     message: err.message,
  //     filename: err.stack,
  //   })
  // }

  // router.onError((event) => {
  //   writeError('Router error', {
  //     name: event.error.name,
  //     message: event.message,
  //     stack: event.stack,
  //   })
  // })

  window.addEventListener('error', function (event) {
    writeError('OnError', {
      name: event.error.name,
      message: event.error.message,
      stack: event.error.stack,
    })
    return true
  })

  window.addEventListener('unhandledrejection', function (event) {
    writeError('unhandledrejection', {
      name: event.reason.name,
      message: event.reason.message,
      stack: event.reason.stack,
    })
    return true
  })

  window.addEventListener('rejectionhandled', function (event) {
    writeError('rejectionhandled', {
      name: event.reason.name,
      message: event.reason.message,
      stack: event.reason.stack,
    })
    return true
  })
}

function $getAllErrors() {
  return out
}

export { $reportErrorGuard, $getAllErrors }
