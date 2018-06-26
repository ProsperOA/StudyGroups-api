Sequel.migration do
  up do
    puts "creating study_groups table"
    create_table(:study_groups) do
      primary_key :id
      foreign_key :user_id,         :users,    :null=>false, :key=>[:id], :on_delete=>:cascade
      String      :name,            :size=>35, :null=>false
      String      :members,         :size=>255
      Integer     :members_limit,   :default=>1
      String      :waitlist,        :size=>255
      Integer     :available_spots, :default=>1
      String      :location,        :size=>140
      String      :description,     :size=>280
      DateTime    :meeting_date
      String      :course
      DateTime    :created_on,      :null=>false
      DateTime    :updated_on,      :null=>false
    end
  end

  down do
    puts "dropping study_groups table"
    drop_table(:study_groups)
  end
end
