$poo

@include('dir2.view2')

@include('ops._nav', [
        'active' => Request::has('requests') ? 'requests' : 'calendar'
])
