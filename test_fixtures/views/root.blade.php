@extends('layouts.master')

@push('content')

  @include('account.sub-header', array(
    'page_title' => $title,
    'back_button' => true,
    ))

<div class="content">

  <div class="row flex flex-start stretch">

    <div class="span12">

  {!! Form::open(array('route' => ['language.store', $locale], 'name' => 'create')) !!}
      <p>{{ ucfirst(trans('app.terminology.explanation')) }}</p>
      <fieldset>
        <div class="fieldset-header">
          <h3>{{ ucfirst(trans_choice('app.terminology.preference', 2)) }}</h3>
        </div>

        @foreach($custom_terms as $termKey => $term)

        <div class="row small">
          <div class="span6 mbs">
            <label for="{{ $termKey }}">{{ $term->getGeneralTerm() }} - {{ trans('app.singular')  }}</label>
            {!! Form::text("terms[{$termKey}][singular]", ucfirst($term->getSingularTerm()), ['id' => $termKey]) !!}
          </div>

          <div class="span6 mbs">
            <label for="{{ $termKey }}_plural">{{ $term->getGeneralTermPlural() }} - {{ trans('app.plural') }}</label>
            {!! Form::text("terms[$termKey][plural]", ucfirst($term->getPluralTerm()), ['id' => $termKey . '_plural']) !!}
          </div>
        </div>

        @endforeach

  </fieldset>

  <div>
    {!! Form::submit(ucfirst(trans('actions.save')), array('class' => 'submit-button')) !!}
    <a href="{{ URL::route('account.index') }}" class="cancel-button">{{ ucfirst(trans('actions.cancel')) }}</a>
  </div>

  {!! Form::close() !!}
</div>
</div>

</div>
@endpush

@push('content-nav')
  @include('account._nav')
@endpush

@push('with-nav') with-nav @endpush
