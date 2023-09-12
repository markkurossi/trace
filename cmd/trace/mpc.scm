;;;
;;; Copyright (c) 2023 Markku Rossi
;;;
;;; All rights reserved.
;;;

;; Trace handler for multi-party computation (MPC) evaluation.
(library (main)
  (export)
  (import (rnrs lists))

  ;; Handle log record.
  (define (handle-record r)
    (let ((message (assoc "message" r))
          (level (assoc "level" r))
          (time (assoc "time" r))
          (attrs (assoc "attrs" r)))
      (display message) (newline)
      (display level) (newline)
      (display time) (newline)
      (display attrs) (newline)
      ))
  )
