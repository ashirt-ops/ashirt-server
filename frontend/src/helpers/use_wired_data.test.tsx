// Copyright 2020, Verizon Media
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

import * as React from 'react'
import { expect } from 'chai'
import { renderHook, act, RenderHookResult } from '@testing-library/react-hooks'
import { spy, SinonSpy } from 'sinon'

import { createDeferred, Deferred } from 'src/test_helpers/deferred'
import { useWiredData, WiredData } from 'src/helpers/use_wired_data'

type DummyData = string

const loadingRenderer = () => <div>loading</div>
const errorRenderer = (err: Error) => <div>error - {err.message}</div>

describe('useWiredData', () => {
  let fetchDataDeferred: Deferred<DummyData>
  let hookTester: RenderHookResult<() => Promise<DummyData>, WiredData<DummyData>>
  let fetchDataSpy: SinonSpy

  beforeEach(() => {
    fetchDataSpy = spy(() => {
      fetchDataDeferred = createDeferred()
      return fetchDataDeferred.promise
    })
    hookTester = renderHook(
      (fetchDataFn: () => Promise<DummyData>) => useWiredData(fetchDataFn, errorRenderer, loadingRenderer),
      { initialProps: fetchDataSpy },
    )
  })

  describe('on initial mount', () => {
    it('calls the fetchDataFn', () => {
      expect(fetchDataSpy).to.have.been.called
    })
    it('reports that it is loading', () => {
      expect(hookTester.result.current.loading).to.equal(true)
    })
    it('renders a loading spinner', () => {
      expect(hookTester.result.current.render(() => <div />)).to.eql(loadingRenderer())
    })
    it('does not call the renderer', () => {
      hookTester.result.current.render(() => { throw Error('Should not be called') })
    })

    describe('reloading while initial data is still loading', () => {
      beforeEach(async () => {
        act(() => { hookTester.result.current.reload() })
      })
      it('does not call the fetchDataFn again', () => {
        expect(fetchDataSpy).to.have.been.calledOnce
      })
    })
  })

  describe('after fetchDataFn returns successfully', () => {
    beforeEach(async () => {
      fetchDataDeferred.resolve('some data')
      await hookTester.waitForNextUpdate()
    })
    it('sets testing to false', () => {
      expect(hookTester.result.current.loading).to.equal(false)
    })
    it('renders the child content by calling renderer with the loaded data', () => {
      const result = hookTester.result.current.render(data => <div>result: {data}</div>)
      expect(result).to.eql(<div>result: {"some data"}</div>)
    })

    describe('reloading after data has finished loading', () => {
      beforeEach(async () => {
        act(() => { hookTester.result.current.reload() })
      })
      it('calls the fetchDataFn', () => {
        expect(fetchDataSpy).to.have.been.calledTwice
      })
      it('displays a loading spinner again', () => {
        expect(hookTester.result.current.render(() => <div />)).to.eql(loadingRenderer())
      })
      it('does not call the renderer', () => {
        hookTester.result.current.render(() => { throw Error('Should not be called') })
      })
    })
  })

  describe('after fetchDataFn returns with an error', () => {
    beforeEach(async () => {
      fetchDataDeferred.reject(Error('some error'))
      await hookTester.waitForNextUpdate()
    })
    it('sets testing to false', () => {
      expect(hookTester.result.current.loading).to.equal(false)
    })
    it('renders an error message', () => {
      expect(hookTester.result.current.render(() => <div />)).to.eql(errorRenderer(Error('some error')))
    })

    describe('reloading after an error', () => {
      beforeEach(async () => {
        act(() => { hookTester.result.current.reload() })
      })
      it('replaces the error with data on success', async () => {
        fetchDataDeferred.resolve('successful data after error')
        await hookTester.waitForNextUpdate()
        const result = hookTester.result.current.render(data => <div>result: {data}</div>)
        expect(result).to.eql(<div>result: {"successful data after error"}</div>)
      })
    })
  })

  describe('debouncing', () => {
    it('debounces rapid changes in fetchDataFn', async () => {
      expect(fetchDataSpy).to.have.callCount(1) // original call from initial mount
      hookTester.rerender(fetchDataSpy.bind({}, '1st rerender'))
      await delay(1)
      hookTester.rerender(fetchDataSpy.bind({}, '2nd rerender'))
      await delay(1)
      hookTester.rerender(fetchDataSpy.bind({}, '3rd rerender'))
      await delay(1)
      hookTester.rerender(fetchDataSpy.bind({}, '4th rerender'))

      // Ensure that we still haven't called it any more since all subsequent rerenders
      // have happened less than 250ms from the initial mount
      expect(fetchDataSpy).to.have.callCount(1)
      expect(fetchDataSpy.getCall(0).args[0]).to.equal(undefined) // initial call in beforeEach did not pass args

      // After 260ms the last rerender should be called:
      await delay(260)
      expect(fetchDataSpy).to.have.callCount(2)
      expect(fetchDataSpy.getCall(1).args[0]).to.equal('4th rerender')

      // After the delay new renders are immediately called:
      hookTester.rerender(fetchDataSpy.bind({}, '5th rerender'))
      expect(fetchDataSpy).to.have.callCount(3)
      expect(fetchDataSpy.getCall(2).args[0]).to.equal('5th rerender')

      // But debouncing is re-enabled:
      hookTester.rerender(fetchDataSpy.bind({}, '6th rerender'))
      expect(fetchDataSpy).to.have.callCount(3)
      expect(fetchDataSpy.getCall(2).args[0]).to.equal('5th rerender')

      // Finally, after another 260ms the last render is once again called:
      await delay(260)
      expect(fetchDataSpy).to.have.callCount(4)
      expect(fetchDataSpy.getCall(3).args[0]).to.equal('6th rerender')
    })

    it('renders debounced data properly', async () => {
      hookTester.rerender(async () => '1st rerender')
      await hookTester.waitForNextUpdate()
      let result = hookTester.result.current.render(data => <div>result: {data}</div>)
      expect(result).to.eql(<div>result: {"1st rerender"}</div>)

      hookTester.rerender(async () => '2nd rerender')
      hookTester.rerender(async () => '3rd rerender')
      hookTester.rerender(async () => '4th rerender')

      // 2nd is rendered immediately since we have been idle
      await hookTester.waitForNextUpdate()
      result = hookTester.result.current.render(data => <div>result: {data}</div>)
      expect(result).to.eql(<div>result: {"2nd rerender"}</div>)

      // 3rd is skipped, and after 250ms 4th is rendered
      await hookTester.waitForNextUpdate()
      result = hookTester.result.current.render(data => <div>result: {data}</div>)
      expect(result).to.eql(<div>result: {"4th rerender"}</div>)
    })
  })
})

const delay = (ms: number) => new Promise(r => setTimeout(r, ms))
