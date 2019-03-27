(defpackage :olduvai
  (:use :cl :postmodern)
  (:export :main))
(in-package :olduvai)

(defun main ()
  (connect-to-database))


(defun connect-to-database ()
  
  (postmodern:connect-toplevel "olduvai" "wrycode" "password" "68.183.55.90"))

;; (postmodern:query "select 22, 'Folie et déraison', 4.5")

(defun infer-type (value)
  )
(defun csv-convert (pathname)
  )

(defun write-db (list?)
  )
(defparameter *ssl-certificate-file* nil)
(defparameter *ssl-key-file* nil)
(defparameter +SSL-VERIFY-NONE+ 1)
