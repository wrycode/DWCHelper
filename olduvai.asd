(asdf:defsystem "olduvai"
  :depends-on (:postmodern)
  :serial t
  :components ((:file "olduvai"))
   :build-operation "program-op" ;; leave as is
 :build-pathname "olduvai"
 :entry-point "olduvai:main")

