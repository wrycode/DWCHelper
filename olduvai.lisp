(defpackage :olduvai
  (:use :cl :postmodern)
  (:export :main))
(in-package :olduvai)

(defun main ()
  (connect-to-database))

(defun connect-to-database ()
  (postmodern:connect-toplevel "mydb" "wrycode" "password" "localhost"))

;; (postmodern:query "select 22, 'Folie et déraison', 4.5")

(defun infer-type (value)
  )
(defun csv-convert (pathname)
  )

(defun write-db (list?)
  )
