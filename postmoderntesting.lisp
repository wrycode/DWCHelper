(asdf:oos 'asdf:load-op :postmodern)
(use-package :postmodern)


(defun main ()
  (connect-toplevel "mydb" "wrycode" "password" "localhost")
  (query "select 22, 'Folie et déraison', 4.5")
  )
