$poo

@include('dir1.dir2.view1')

@include('ops._nav', [
        'active' => Request::has('requests') ? 'requests' : 'calendar'
])
