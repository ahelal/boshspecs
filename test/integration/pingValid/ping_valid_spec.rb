describe command('boshspecs ping') do
    its('exit_status') { should eq 0 }
    its('stdout') { should match (/bbl-bosh-gcp/) }
    its('stdout') { should match (/OK/) }
  end