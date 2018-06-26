Sequel.migration do
  up do
    puts "creating users table"
    create_table(:users, :ignore_index_errors=>true) do
      primary_key :id
      String      :first_name,   :size=>20,  :null=>false
      String      :last_name,    :size=>20
      String      :email,        :size=>60,  :null=>false
      String      :password,     :size=>80,  :null=>false
      String      :avatar,       :size=>255, :default=>"https://s3-us-west-1.amazonaws.com/study-groups/images/user-avatars/stock-avatar.png"
      String      :bio,          :size=>280
      String      :school,       :size=>20
      String      :major1,       :size=>40
      String      :major2,       :size=>40
      String      :minor,        :size=>40
      String      :study_groups, :size=>255
      String      :waitlists,    :size=>255
      String      :courses
      DateTime    :created_on,   :null=>false
      DateTime    :updated_on,   :null=>false

      index [:email], :name=>:users_email_key, :unique=>true
    end
  end

  down do
    puts "dropping users table"
    drop_table(:users, :cascade => true)
  end
end
