From Coq Require Import List String Ascii Bool.
Import ListNotations.

Open Scope string_scope.

Fixpoint split_semi_aux (acc : list ascii) (chars : list ascii) : list string :=
  match chars with
  | nil => [string_of_list_ascii acc]
  | cons c rest =>
      if Ascii.eqb c ";"%char
      then string_of_list_ascii acc :: split_semi_aux nil rest
      else split_semi_aux (acc ++ [c]) rest
  end.

Definition split_semi (s : string) : list string :=
  split_semi_aux nil (list_ascii_of_string s).

Definition trim (s : string) : string := s.

Definition not_empty (s : string) : bool :=
  if String.eqb s "" then false else true.

Definition normalize_stmt (stmt : string) : list string :=
  filter (fun q => not_empty (trim q)) (split_semi stmt).

Fixpoint concat_strings (l : list (list string)) : list string :=
  match l with
  | nil => nil
  | cons h t => h ++ concat_strings t
  end.

Definition normalize (commands : list string) : list string :=
  concat_strings (map normalize_stmt commands).
